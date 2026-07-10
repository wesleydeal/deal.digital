package main

import "testing"

func TestParseCLIAcceptsFlagsBeforeAndAfterCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		command   string
		drafts    bool
		addr      string
		outputDir string
	}{
		{name: "flag before command", args: []string{"--drafts", "build"}, command: "build", drafts: true},
		{name: "flag after command", args: []string{"build", "--drafts"}, command: "build", drafts: true},
		{name: "value flag after command", args: []string{"serve", "--addr", ":9090"}, command: "serve", addr: ":9090"},
		{name: "value flag before command", args: []string{"--output", "dist", "build"}, command: "build", outputDir: "dist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := parseCLI(tt.args)
			if err != nil {
				t.Fatalf("parse CLI: %v", err)
			}
			if cli.command != tt.command {
				t.Fatalf("command = %q, want %q", cli.command, tt.command)
			}
			if cli.opts.IncludeDrafts != tt.drafts {
				t.Fatalf("drafts = %v, want %v", cli.opts.IncludeDrafts, tt.drafts)
			}
			if tt.addr != "" && cli.addr != tt.addr {
				t.Fatalf("addr = %q, want %q", cli.addr, tt.addr)
			}
			if tt.outputDir != "" && cli.opts.OutputDir != tt.outputDir {
				t.Fatalf("output = %q, want %q", cli.opts.OutputDir, tt.outputDir)
			}
		})
	}
}

func TestParseCLIRejectsSingleDashLongFlags(t *testing.T) {
	if _, err := parseCLI([]string{"build", "-drafts"}); err == nil {
		t.Fatal("single-dash long flag was accepted")
	}
}
