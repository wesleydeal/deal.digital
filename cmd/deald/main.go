package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"deal.digital/internal/content"
	"deal.digital/internal/site"
)

func main() {
	log.SetFlags(0)

	addr := flag.String("addr", ":8080", "listen address")
	configPath := flag.String("config", "config.toml", "path to site config")
	contentDir := flag.String("content", "content", "content directory")
	staticDir := flag.String("static", "static", "static directory")
	templateDir := flag.String("templates", "templates", "template directory")
	outputDir := flag.String("output", filepath.Join("build", "site"), "rendered site output directory")
	includeDrafts := flag.Bool("drafts", false, "include draft content")
	flag.Parse()

	command := "serve"
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	opts := content.Options{
		ConfigPath:    *configPath,
		ContentDir:    *contentDir,
		StaticDir:     *staticDir,
		TemplateDir:   *templateDir,
		OutputDir:     *outputDir,
		IncludeDrafts: *includeDrafts,
	}
	if command == "build" {
		opts.CompressOutput = true
	}

	app, err := site.NewGenerator(opts)
	if err != nil {
		log.Fatal(err)
	}

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
		app.Watch(ctx)
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
