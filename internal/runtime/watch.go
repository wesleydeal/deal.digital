package runtime

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"
)

func Watch(ctx context.Context, roots []string, onChange func(context.Context) error) {
	if len(roots) == 0 {
		return
	}

	go func() {
		snapshot, err := snapshotFiles(roots)
		if err != nil {
			log.Printf("watch disabled: %v", err)
			return
		}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				next, err := snapshotFiles(roots)
				if err != nil {
					log.Printf("watch scan failed: %v", err)
					continue
				}
				if snapshotsEqual(snapshot, next) {
					continue
				}

				snapshot = next
				log.Printf("changes detected, rebuilding site")
				if err := onChange(ctx); err != nil {
					log.Printf("rebuild failed: %v", err)
				}
			}
		}
	}()
}

func NormalizeRoots(roots []string) []string {
	seen := map[string]struct{}{}
	unique := make([]string, 0, len(roots))
	for _, root := range roots {
		if root == "" {
			continue
		}
		clean := filepath.Clean(root)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		unique = append(unique, clean)
	}
	slices.Sort(unique)
	return unique
}

type fileStamp struct {
	Size    int64
	ModTime time.Time
}

type watchSnapshot map[string]fileStamp

func snapshotFiles(roots []string) (watchSnapshot, error) {
	snapshot := watchSnapshot{}
	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if !info.IsDir() {
			snapshot[root] = fileStamp{Size: info.Size(), ModTime: info.ModTime().UTC()}
			continue
		}
		if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			snapshot[path] = fileStamp{Size: info.Size(), ModTime: info.ModTime().UTC()}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return snapshot, nil
}

func snapshotsEqual(a, b watchSnapshot) bool {
	if len(a) != len(b) {
		return false
	}
	for path, stamp := range a {
		other, ok := b[path]
		if !ok || stamp != other {
			return false
		}
	}
	return true
}
