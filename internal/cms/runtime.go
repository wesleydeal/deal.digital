package cms

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func compressTree(root string) error {
	files := []string{}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		switch filepath.Ext(path) {
		case ".gz", ".br", ".zst":
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return err
	}

	for _, path := range files {
		plain, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := writeCompressedVariants(path, plain); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) ListenAndServe(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/_/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/_/admin/rebuild", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := a.Rebuild(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	mux.HandleFunc("/_/fragments/sample-counter", a.handleSampleCounter)

	fsHandler := http.FileServer(http.Dir(a.opts.OutputDir))
	mux.Handle("/", fsHandler)

	server := &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(mux),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/_/") {
			w.Header().Set("Cache-Control", "public, max-age=60")
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) handleSampleCounter(w http.ResponseWriter, r *http.Request) {
	var (
		value int64
		err   error
	)

	switch r.Method {
	case http.MethodGet:
		value, err = a.store.CounterValue(r.Context(), "sample-page")
	case http.MethodPost:
		value, err = a.store.IncrementCounter(r.Context(), "sample-page")
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(sampleCounterFragment(value)))
}

func sampleCounterFragment(value int64) string {
	return fmt.Sprintf(`<p>SQLite-backed sample counter: <output aria-live="polite">%d</output></p>`, value)
}
