# gr

gr is a task runner focused on simplicity and minimalism.

It's heavily inspired by just the scrpits section of package.json. To get started create a gr.toml (or run `gr init`), and add your commands, then pass them as an argument to run them.

For example the simplest file

```toml
[commands]
hello = "echo Hello, world"
```

to run it simply call `gr hello`.

## Config Fields

This section decribes the different attributes that can be added to the config in the gr.toml file.

`[commands]`

The commands that are able to be run. They can be set in two varieties, simple commands and detailed commands. Simple commands are entered just as you would type them on the command line. Detailed commands are created using toml's tables.

```toml
[commands]
build = "go build -o gr ."
hello = "echo hello, world"
goodbye = "echo goodbye, world"
both = ["hello", "goodbye"]

[commands.clean]
command = "rm"
args = ["-rf", "./node_modules"]
```

`[environment]`

A table of all environment variables that you want set on every command.

```toml
[environment]
TARGET = "world"

[commands]
greet = "hello $TARGET"
```

`working_directory`

Set the directory from which all commands run.

### Detailed Commands

This section details the settings available to detailed commands. On any setting applied to both the global config, and an individual command, the command will take priority.

```toml
[commands.list-deps]
command = "ls"
args = ["-a", "."]
working_directory = "./node_modules"
environment = { FOO="bar" }
```

`command` : The command to run.
`args` : The list of strings to pass as arguments to the command
`working_directory` : The directory from which to run the command
`environment` : A table of environment variables to add to the command
`steps` : The commands to run if you are doing a multi-stage command.

## CLI

The Cli takes the following form. `gr [options] <command>`, with the singular command being required, and multiple optional options.

### Commands

- `gr help` Prints out the help text.
- `gr list` Prints all of the available commands that you can run.
- `gr init` Creates a new `gr.toml` file.
- `gr <command>` Runs the command specified in the config file.

### Options

- `-f, --file` Specifile a filename other than gr.toml for the config
- `--logs <file>` Print output to the log file.
- `-q, --quiet` Silences all output from the terminal.
