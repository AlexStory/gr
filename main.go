package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type options struct {
	quiet   bool
	logs    string
	command string
}

func main() {
	var quiet bool
	var logs string

	flag.BoolVar(&quiet, "q", false, "suppress output")
	flag.BoolVar(&quiet, "quiet", false, "suppress output")
	flag.StringVar(&logs, "logs", "", "write output to file")

	flag.Usage = helpCmd

	flag.Parse()
	cmd := flag.Arg(0)
	opts := options{quiet: quiet, logs: logs, command: cmd}

	switch cmd {
	case "list":
		listCmd()
	case "", "help":
		helpCmd()
	default:
		runCmd(opts)
	}
}

func helpCmd() {
	fmt.Println("Usage: gr [options] <command>")
	fmt.Println("Commands:")
	fmt.Println("  help  prints this message")
	fmt.Println("  list  lists all available commans")
	fmt.Println("\nOptions:")
	fmt.Println("  -q, --quiet    suppress output")
	fmt.Println("  --logs <file>  write output to file")
}

func listCmd() {
	config := loadConfig()

	for _, task := range config.Commands {
		fmt.Println(task.Name)
	}
}

func runCmd(opts options) {
	config := loadConfig()

	for _, task := range config.Commands {
		if task.Name == opts.command {
			runTask(task, opts)
			return
		}
	}

	fmt.Printf("Command %q not found\n", opts.command)
}

func runTask(task Command, opts options) {
	cmd := exec.Command(task.Command, task.Arguments...)
	var outputWriter io.Writer

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
			outputWriter = io.MultiWriter(file, os.Stdout)
		}
	} else {
		if opts.quiet {
			outputWriter = io.Discard
		} else {
			outputWriter = os.Stdout
		}
	}

	cmd.Stdout = outputWriter
	cmd.Stderr = os.Stderr
	cmd.Run()
}
