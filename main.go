package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const version = "0.1.1"

type options struct {
	quiet      bool
	logs       string
	command    string
	configFile string
	writer     io.Writer
}

func main() {
	opts := parseFlags()
	opts.writer = os.Stdout

	switch opts.command {
	case "list":
		listCmd(opts)
	case "", "help":
		helpCmd(opts)
	case "init":
		initCmd(opts)
	case "version":
		versionCmd(opts)
	default:
		runCmd(opts)
	}
}

func parseFlags() *options {
	var quiet bool
	var configFile string
	var logs string

	flag.BoolVar(&quiet, "q", false, "suppress output")
	flag.BoolVar(&quiet, "quiet", false, "suppress output")
	flag.StringVar(&logs, "logs", "", "write output to file")
	flag.StringVar(&configFile, "f", "gr.toml", "config file")
	flag.StringVar(&configFile, "file", "gr.toml", "config file")

	flag.Parse()
	cmd := flag.Arg(0)
	opts := &options{
		quiet:      quiet,
		logs:       logs,
		command:    cmd,
		configFile: configFile,
	}
	flag.Usage = func() { helpCmd(opts) }
	return opts
}

func helpCmd(opts *options) {
	fmt.Fprintln(opts.writer, "Usage: gr [options] <command>")
	fmt.Fprintln(opts.writer, "Commands:")
	fmt.Fprintln(opts.writer, "  help     prints this message")
	fmt.Fprintln(opts.writer, "  init     creates a new config file")
	fmt.Fprintln(opts.writer, "  list     lists all available commans")
	fmt.Fprintln(opts.writer, "  version  prints the version")
	fmt.Fprintln(opts.writer, "\nOptions:")
	fmt.Fprintln(opts.writer, "  -f, --file <file> specify the config file (default: gr.toml)")
	fmt.Fprintln(opts.writer, "  --logs <file>     write output to file")
	fmt.Fprintln(opts.writer, "  -q, --quiet       suppress output")
}

func listCmd(opts *options) {
	config := loadConfig(opts.configFile)

	fmt.Fprintln(opts.writer, "Available commands:")

	for _, task := range config.Commands {
		fmt.Fprintf(opts.writer, " - %s\n", task.Name)
	}
}

func initCmd(opts *options) {
	if _, err := os.Stat(opts.configFile); err == nil {
		fmt.Println("Config file already exists")
		return
	}

	file, err := os.Create(opts.configFile)
	if err != nil {
		fmt.Println("Failed to create config file:", err)
		return
	}
	defer file.Close()

	defaultConfig := `[commands]
hello = "echo Hello, World!"`

	file.WriteString(defaultConfig)
}

func versionCmd(opts *options) {
	fmt.Fprint(opts.writer, version)
}

func runCmd(opts *options) {
	config := loadConfig(opts.configFile)

	for _, task := range config.Commands {
		if task.Name == opts.command {
			runTask(task, config, opts)
			return
		}
	}

	fmt.Fprintf(opts.writer, "Command %q not found\n", opts.command)
}

func runTask(task Command, config *Config, opts *options) {
	if task.Steps != nil {
		for _, step := range task.Steps {
			o := &options{
				command:    step,
				quiet:      opts.quiet,
				logs:       opts.logs,
				configFile: opts.configFile,
				writer:     opts.writer,
			}
			runCmd(o)
		}
		return
	}

	cmd := exec.Command(task.Command, task.Arguments...)
	cmd.Env = append(os.Environ(), formatEnv(config.Environment)...)
	cmd.Env = append(cmd.Env, formatEnv(task.Environment)...)
	var outputWriter io.Writer

	if config.WorkingDirectory != "" {
		cmd.Dir = config.WorkingDirectory
	}

	if task.WorkingDirectory != "" {
		cmd.Dir = task.WorkingDirectory
	}

	if opts.logs != "" {
		file, err := os.OpenFile(opts.logs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("failed to open log file:", err)
			return
		}
		defer file.Close()

		if opts.quiet {
			outputWriter = file
		} else {
			outputWriter = io.MultiWriter(file, opts.writer)
		}
	} else {
		if opts.quiet {
			outputWriter = io.Discard
		} else {
			outputWriter = opts.writer
		}
	}

	cmd.Stdout = outputWriter
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func formatEnv(env map[string]string) []string {
	var formatted []string
	for key, value := range env {
		formatted = append(formatted, fmt.Sprintf("%s=%s", key, value))
	}
	return formatted
}
