package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type options struct {
	quiet bool
}

func main() {
	var quiet bool

	flag.BoolVar(&quiet, "q", false, "suppress output")
	flag.BoolVar(&quiet, "quiet", false, "suppress output")

	flag.Usage = helpCmd

	flag.Parse()
	opts := options{quiet: quiet}
	cmd := flag.Arg(0)

	switch cmd {
	case "list":
		listCmd()
	case "", "help":
		helpCmd()
	default:
		runCmd(cmd, opts)
	}
}

func helpCmd() {
	fmt.Println("Usage: gr [options] <command>")
	fmt.Println("Commands:")
	fmt.Println("  help  prints this message")
	fmt.Println("  list  lists all available commans")
	fmt.Println("\nOptions:")
	fmt.Println("  -q, --quiet  suppress output")
}

func listCmd() {
	config := loadConfig()

	for _, task := range config.Commands {
		fmt.Println(task.Name)
	}
}

func runCmd(cmd string, opts options) {
	config := loadConfig()

	for _, task := range config.Commands {
		if task.Name == cmd {
			runTask(task, opts)
			return
		}
	}

	fmt.Printf("Command %q not found\n", cmd)
}

func runTask(task Command, opts options) {
	cmd := exec.Command(task.Command, task.Arguments...)
	if !opts.quiet {
		fmt.Printf("Running %s %v...\n\n", task.Command, strings.Join(task.Arguments, " "))
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = nil
	}
	cmd.Stderr = os.Stderr
	cmd.Run()
}
