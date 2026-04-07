package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"deal.digital/internal/site"
)

func main() {
	log.SetFlags(0)

	addr := flag.String("addr", ":8080", "listen address")
	configPath := flag.String("config", "config.toml", "path to site config")
	contentDir := flag.String("content", "content", "content directory")
	staticDir := flag.String("static", "static", "static directory")
	outputDir := flag.String("output", filepath.Join("build", "site"), "rendered site output directory")
	dataDir := flag.String("data", filepath.Join("build", "data"), "runtime data directory")
	includeDrafts := flag.Bool("drafts", false, "include draft content")
	flag.Parse()

	command := "serve"
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	app, err := site.NewApp(site.Options{
		ConfigPath:    *configPath,
		ContentDir:    *contentDir,
		StaticDir:     *staticDir,
		OutputDir:     *outputDir,
		DataDir:       *dataDir,
		IncludeDrafts: *includeDrafts,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	ctx := signalContext()
	switch command {
	case "build":
		if err := app.Rebuild(ctx); err != nil {
			log.Fatal(err)
		}
	case "serve":
		if err := app.Rebuild(ctx); err != nil {
			log.Fatal(err)
		}
		log.Printf("serving %s on %s", app.Options().OutputDir, *addr)
		if err := app.ListenAndServe(ctx, *addr); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command %q (expected build or serve)", command)
	}
}

func signalContext() context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return ctx
}
