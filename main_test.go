package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
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
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("tt #%d", i), func(t *testing.T) {
			listCmd(tt.opts)
			if !strings.Contains(tt.opts.writer.(*bytes.Buffer).String(), "- hello") {
				t.Errorf("listCmd() = %v, want %v", tt.opts.writer.(*bytes.Buffer).String(), tt.expected)
			}

			if !strings.Contains(tt.opts.writer.(*bytes.Buffer).String(), "- goodbye") {
				t.Errorf("listCmd() = %v, want %v", tt.opts.writer.(*bytes.Buffer).String(), tt.expected)
			}
		})
	}
}

func TestFormatEnv(t *testing.T) {
	tests := []struct {
		env  map[string]string
		want []string
	}{
		{
			env: map[string]string{
				"key1": "value1",
			},
			want: []string{"key1=value1"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("tt #%d", i), func(t *testing.T) {
			got := formatEnv(tt.env)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelpCmd(t *testing.T) {
	buffer := new(bytes.Buffer)
	opts := &options{
		writer: buffer,
	}

	helpCmd(opts)

	got := buffer.String()
	want := `Usage: gr [options] <command>
Commands:
  help  prints this message
  init  creates a new config file
  list  lists all available commans

Options:
  -f, --file <file> specify the config file (default: gr.toml)
  --logs <file>     write output to file
  -q, --quiet       suppress output
`

	if got != want {
		t.Errorf("helpCmd() = %q, want %q", got, want)
	}
}

func TestRunCmd(t *testing.T) {
	tests := []struct {
		opts     *options
		expected string
	}{
		{
			opts: &options{
				quiet:      true,
				logs:       "",
				command:    "doesntexist",
				configFile: "test.toml",
				writer:     new(bytes.Buffer),
			},
			expected: "Command \"doesntexist\" not found\n",
		},
		{
			opts: &options{
				quiet:      false,
				logs:       "",
				writer:     new(bytes.Buffer),
				command:    "hello",
				configFile: "test.toml",
			},
			expected: "hello, world\n",
		},
		{
			opts: &options{
				quiet:      false,
				logs:       "",
				writer:     new(bytes.Buffer),
				command:    "both",
				configFile: "test.toml",
			},
			expected: "hello, world\ngoodbye, world\n",
		},
		{
			opts: &options{
				quiet:      true,
				logs:       "",
				writer:     new(bytes.Buffer),
				command:    "both",
				configFile: "test.toml",
			},
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("tt #%d", i), func(t *testing.T) {
			runCmd(tt.opts)
			if tt.opts.writer.(*bytes.Buffer).String() != tt.expected {
				t.Errorf("runCmd() = %v, want %v", tt.opts.writer.(*bytes.Buffer).String(), tt.expected)
			}
		})
	}
}

func TestVersionCmd(t *testing.T) {
	buffer := new(bytes.Buffer)
	opts := &options{
		writer: buffer,
	}

	versionCmd(opts)

	got := buffer.String()
	want := version

	if got != want {
		t.Errorf("versionCmd() = %q, want %q", got, want)
	}
}
