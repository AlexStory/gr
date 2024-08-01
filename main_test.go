package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		args        []string
		wantOptions *options
	}{
		{
			args: []string{"-q", "list"},
			wantOptions: &options{
				quiet:      true,
				logs:       "",
				command:    "list",
				configFile: "gr.toml",
			},
		},
		{
			args: []string{"-logs", "output.log", "init"},
			wantOptions: &options{
				quiet:      false,
				logs:       "output.log",
				command:    "init",
				configFile: "gr.toml",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("tt #%d", i), func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = append([]string{"cmd"}, tt.args...)
			gotOptions := parseFlags()
			if *gotOptions != *tt.wantOptions {
				t.Errorf("parseFlags() = %v, want %v", gotOptions, tt.wantOptions)
			}
		})
	}
}

func TestListCmd(t *testing.T) {
	tests := []struct {
		opts     *options
		expected string
	}{
		{
			opts: &options{
				quiet:      true,
				logs:       "",
				command:    "list",
				configFile: "test.toml",
				writer:     new(bytes.Buffer),
			},
			expected: `Available commands:
 - hello
 - goodbye
`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("tt #%d", i), func(t *testing.T) {
			listCmd(tt.opts)
			if tt.opts.writer.(*bytes.Buffer).String() != tt.expected {
				t.Errorf("listCmd() = %q, want %q", tt.opts.writer.(*bytes.Buffer).String(), tt.expected)
			}
		})
	}
}
