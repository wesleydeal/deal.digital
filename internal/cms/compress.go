package cms

import (
	"bytes"
	"compress/gzip"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

type compressionPlan struct {
	gzip   bool
	brotli bool
	zstd   bool
}

type compressionJob struct {
	path string
	plan compressionPlan
}

var compressionBlacklist = map[string]struct{}{
	".7z":    {},
	".aac":   {},
	".avif":  {},
	".br":    {},
	".bz2":   {},
	".dll":   {},
	".dmg":   {},
	".exe":   {},
	".gif":   {},
	".gz":    {},
	".ico":   {},
	".iso":   {},
	".jpeg":  {},
	".jpg":   {},
	".lz4":   {},
	".m4a":   {},
	".mov":   {},
	".mp3":   {},
	".mp4":   {},
	".msi":   {},
	".ogg":   {},
	".oga":   {},
	".ogv":   {},
	".png":   {},
	".pdf":   {},
	".rar":   {},
	".swf":   {},
	".webm":  {},
	".webp":  {},
	".woff2": {},
	".xz":    {},
	".zip":   {},
	".zst":   {},
}

var brotliWhitelist = map[string]struct{}{
	".atom":        {},
	".css":         {},
	".htm":         {},
	".html":        {},
	".ics":         {},
	".js":          {},
	".json":        {},
	".map":         {},
	".md":          {},
	".mjs":         {},
	".rss":         {},
	".svg":         {},
	".txt":         {},
	".webmanifest": {},
	".xml":         {},
	".xhtml":       {},
	".yml":         {},
	".yaml":        {},
}

func compressTree(root string) error {
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}

	jobs := make(chan compressionJob, workers*2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	var errOnce sync.Once
	recordErr := func(err error) {
		if err == nil {
			return
		}
		errOnce.Do(func() {
			errCh <- err
			cancel()
		})
	}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					if err := compressFile(job.path, job.plan); err != nil {
						recordErr(err)
						return
					}
				}
			}
		}()
	}

	walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		plan, ok := compressionPlanFor(path)
		if !ok {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case jobs <- compressionJob{path: path, plan: plan}:
			return nil
		}
	})

	close(jobs)
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return walkErr
	}
}

func compressionPlanFor(path string) (compressionPlan, bool) {
	ext := strings.ToLower(filepath.Ext(path))
	if _, blacklisted := compressionBlacklist[ext]; blacklisted {
		return compressionPlan{}, false
	}

	plan := compressionPlan{
		gzip: true,
		zstd: true,
	}
	if _, allowed := brotliWhitelist[ext]; allowed {
		plan.brotli = true
	}
	return plan, true
}

func compressFile(path string, plan compressionPlan) error {
	plain, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return writeCompressedVariants(path, plain, plan)
}

func writeCompressedVariants(basePath string, plain []byte, plan compressionPlan) error {
	if plan.gzip {
		if err := writeIfSmaller(basePath+".gz", plain, gzipBytes); err != nil {
			return err
		}
	}
	if plan.brotli {
		if err := writeIfSmaller(basePath+".br", plain, brotliBytes); err != nil {
			return err
		}
	}
	if plan.zstd {
		if err := writeIfSmaller(basePath+".zst", plain, zstdBytes); err != nil {
			return err
		}
	}
	return nil
}

func writeIfSmaller(path string, plain []byte, encode func([]byte) ([]byte, error)) error {
	compressed, err := encode(plain)
	if err != nil {
		return err
	}
	if len(compressed) >= len(plain) {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	return os.WriteFile(path, compressed, 0o644)
}

func gzipBytes(plain []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := zw.Write(plain); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func brotliBytes(plain []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := brotli.NewWriterLevel(&buf, brotli.BestCompression)
	if _, err := zw.Write(plain); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func zstdBytes(plain []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw, err := zstd.NewWriter(&buf, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return nil, err
	}
	if _, err := zw.Write(plain); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
