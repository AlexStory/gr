package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
)

var re = regexp.MustCompile(`'[^']*'|"[^"]*"|\S+`)

type Command struct {
	Name             string
	Command          string
	Arguments        []string
	Steps            []string
	Environment      map[string]string
	WorkingDirectory string
}

type Config struct {
	Commands         []Command
	Environment      map[string]string
	WorkingDirectory string
}

func parseConfig(data map[string]interface{}) *Config {
	config := &Config{
		Environment: make(map[string]string),
	}

	if commands, ok := data["commands"].(map[string]interface{}); ok {
		for name, cmd := range commands {
			switch cmd := cmd.(type) {
			case string:
				command := splitCommand(cmd)[0]
				args := splitCommand(cmd)[1:]
				config.Commands = append(config.Commands, Command{
					Name:      name,
					Command:   command,
					Arguments: args,
				})
			case []interface{}:
				command := Command{Name: name}
				for _, step := range cmd {
					if stepStr, ok := step.(string); ok {
						command.Steps = append(command.Steps, stepStr)
					}
				}
				config.Commands = append(config.Commands, command)

			case map[string]interface{}:
				command := Command{Name: name}
				if cmdStr, ok := cmd["command"].(string); ok {
					command.Command = cmdStr
				}
				if args, ok := cmd["args"].([]interface{}); ok {
					for _, arg := range args {
						if argStr, ok := arg.(string); ok {
							command.Arguments = append(command.Arguments, argStr)
						}
					}
				}
				if steps, ok := cmd["steps"].([]interface{}); ok {
					for _, step := range steps {
						if stepStr, ok := step.(string); ok {
							command.Steps = append(command.Steps, stepStr)
						}
					}
				}

				if env, ok := cmd["environment"].(map[string]interface{}); ok {
					command.Environment = make(map[string]string)
					for key, value := range env {
						if valueStr, ok := value.(string); ok {
							command.Environment[key] = valueStr
						}
					}
				}

				if wd, ok := cmd["working_directory"].(string); ok {
					command.WorkingDirectory = wd
				}

				config.Commands = append(config.Commands, command)
			}
		}
	}

	if env, ok := data["environment"].(map[string]interface{}); ok {
		for key, value := range env {
			if valueStr, ok := value.(string); ok {
				config.Environment[key] = valueStr
			}
		}
	}

	if wd, ok := data["working_directory"].(string); ok {
		config.WorkingDirectory = wd
	}

	return config
}

func loadConfig(filename string) *Config {
	configPath, err := findConfigFile(filename)
	if err != nil {
		fmt.Println("config file not found")
		os.Exit(1)
	}

	var data map[string]interface{}
	if _, err := toml.DecodeFile(configPath, &data); err != nil {
		fmt.Println("failed to parse config file:", err)
		os.Exit(1)
	}

	return parseConfig(data)
}

func findConfigFile(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		configPath := filepath.Join(dir, filename)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("config file not found")
}

func splitCommand(cmd string) []string {
	return re.FindAllString(cmd, -1)
}
