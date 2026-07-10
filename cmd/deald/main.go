package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"deal.digital/internal/content"
	"deal.digital/internal/site"
)

func main() {
	log.SetFlags(0)
	cli, err := parseCLI(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	app, err := site.NewGenerator(cli.opts)
	if err != nil {
		log.Fatal(err)
	}

	ctx := signalContext()
	switch cli.command {
	case "build":
		if err := app.Rebuild(ctx); err != nil {
			log.Fatal(err)
		}
	case "serve":
		if err := app.Rebuild(ctx); err != nil {
			log.Fatal(err)
		}
		app.Watch(ctx)
		log.Printf("serving %s on %s", app.Options().OutputDir, cli.addr)
		if err := app.ListenAndServe(ctx, cli.addr); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command %q (expected build or serve)", cli.command)
	}
}

type cliConfig struct {
	command string
	addr    string
	opts    content.Options
}

func parseCLI(args []string) (cliConfig, error) {
	command, flagArgs, err := splitCommandAndFlags(args)
	if err != nil {
		return cliConfig{}, err
	}

	flags := flag.NewFlagSet("deald", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: deald [build|serve] [flags]")
		fmt.Fprintln(flags.Output(), "")
		fmt.Fprintln(flags.Output(), "Flags:")
		fmt.Fprintln(flags.Output(), "  --addr <address>       listen address for serve (default :8080)")
		fmt.Fprintln(flags.Output(), "  --config <path>        site configuration (default config.toml)")
		fmt.Fprintln(flags.Output(), "  --content <path>       content directory (default content)")
		fmt.Fprintln(flags.Output(), "  --static <path>        static asset directory (default static)")
		fmt.Fprintln(flags.Output(), "  --templates <path>     template directory (default templates)")
		fmt.Fprintln(flags.Output(), "  --output <path>        output directory (default build/site)")
		fmt.Fprintln(flags.Output(), "  --drafts               include draft content")
	}

	addr := flags.String("addr", ":8080", "")
	configPath := flags.String("config", "config.toml", "")
	contentDir := flags.String("content", "content", "")
	staticDir := flags.String("static", "static", "")
	templateDir := flags.String("templates", "templates", "")
	outputDir := flags.String("output", filepath.Join("build", "site"), "")
	includeDrafts := flags.Bool("drafts", false, "")
	if err := flags.Parse(flagArgs); err != nil {
		return cliConfig{}, err
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
	return cliConfig{command: command, addr: *addr, opts: opts}, nil
}

var longFlags = map[string]bool{
	"addr":      true,
	"config":    true,
	"content":   true,
	"static":    true,
	"templates": true,
	"output":    true,
	"drafts":    false,
	"help":      false,
}

func splitCommandAndFlags(args []string) (string, []string, error) {
	command := "serve"
	commandSet := false
	flagArgs := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if !strings.HasPrefix(arg, "--") || arg == "--" {
				return "", nil, fmt.Errorf("use -- for flags: %q", arg)
			}
			name, _, hasValue := strings.Cut(strings.TrimPrefix(arg, "--"), "=")
			needsValue, ok := longFlags[name]
			if !ok {
				return "", nil, fmt.Errorf("unknown flag %q", arg)
			}
			flagArgs = append(flagArgs, arg)
			if needsValue && !hasValue {
				if i+1 == len(args) {
					return "", nil, fmt.Errorf("flag %q requires a value", arg)
				}
				i++
				flagArgs = append(flagArgs, args[i])
			}
			continue
		}
		if commandSet {
			return "", nil, fmt.Errorf("unexpected argument %q", arg)
		}
		if arg != "build" && arg != "serve" {
			return "", nil, fmt.Errorf("unknown command %q (expected build or serve)", arg)
		}
		command = arg
		commandSet = true
	}
	return command, flagArgs, nil
}

func signalContext() context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return ctx
}
