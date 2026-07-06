package site

import (
	"context"
	"net/http"
	"time"
)

func (a *Generator) ListenAndServe(ctx context.Context, addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: http.FileServer(http.Dir(a.opts.OutputDir)),
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
