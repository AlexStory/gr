package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	t.Run("test detailed", func(t *testing.T) {
		got := loadConfig("./testFiles/detailed.toml")
		want := &Config{
			Commands: []Command{
				{
					Name:             "detailed",
					Command:          "echo",
					Arguments:        []string{"hello, world"},
					Steps:            []string{},
					Environment:      map[string]string{},
					WorkingDirectory: "",
				},
			},
			Environment:      map[string]string{},
			WorkingDirectory: "",
		}

		if got.Commands[0].Name != (want.Commands[0].Name) {
			t.Errorf("got %v, want %v", got.Commands[0], want.Commands[0])
		}

		if got.Commands[0].Command != (want.Commands[0].Command) {
			t.Errorf("got %v, want %v", got.Commands[0], want.Commands[0])
		}

		if !reflect.DeepEqual(got.Commands[0].Arguments, want.Commands[0].Arguments) {
			t.Errorf("got %v, want %v", got.Commands[0], want.Commands[0])
		}
	})

	t.Run("test simple steps", func(t *testing.T) {
		got := loadConfig("./testFiles/steps.toml")
		want := &Config{
			Commands: []Command{
				{
					Name:             "both",
					Command:          "",
					Arguments:        []string{},
					Steps:            []string{"hello", "goodbye"},
					Environment:      map[string]string{},
					WorkingDirectory: "",
				},
			},
			Environment:      map[string]string{},
			WorkingDirectory: "",
		}

		var res Command
		for _, c := range got.Commands {
			if len(c.Steps) > 0 {
				res = c
			}
		}

		if res.Name != (want.Commands[0].Name) {
			t.Errorf("got %v, want %v", res.Name, want.Commands[0].Name)
		}

		if !reflect.DeepEqual(res.Steps, want.Commands[0].Steps) {
			t.Errorf("got %v, want %v", res, want.Commands[0])
		}
	})

	t.Run("test detailed options", func(t *testing.T) {
		got := loadConfig("./testFiles/options.toml")
		want := &Config{
			Commands: []Command{
				{
					Name:             "detailed",
					Command:          "echo",
					Arguments:        []string{"hello, world"},
					Steps:            []string{},
					Environment:      map[string]string{"FOO": "bar"},
					WorkingDirectory: "./out",
				},
			},
		}

		if got.Commands[0].Name != (want.Commands[0].Name) {
			t.Errorf("Name: got %v, want %v", got.Commands[0], want.Commands[0])
		}

		if got.Commands[0].Command != (want.Commands[0].Command) {
			t.Errorf("Command: got %v, want %v", got.Commands[0], want.Commands[0])
		}

		if !reflect.DeepEqual(got.Commands[0].Arguments, want.Commands[0].Arguments) {
			t.Errorf("Arguments: got %v, want %v", got.Commands[0].Arguments, want.Commands[0].Arguments)
		}

		if !reflect.DeepEqual(got.Commands[0].Environment, want.Commands[0].Environment) {
			t.Errorf("ENV: got %v, want %v", got.Commands[0], want.Commands[0])
		}

		if got.Commands[0].WorkingDirectory != (want.Commands[0].WorkingDirectory) {
			t.Errorf("Working_dir: got %v, want %v", got.Commands[0], want.Commands[0])
		}
	})
}

func TestFindConfigFile(t *testing.T) {
	t.Run("test find config", func(t *testing.T) {
		got, _ := findConfigFile("gr.toml")
		want := "gr.toml"

		dir, _ := os.Getwd()
		relativePath, _ := filepath.Rel(dir, got)

		if relativePath != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("return error when not found", func(t *testing.T) {
		_, err := findConfigFile("notfound.toml")

		if err == nil {
			t.Errorf("expected an error, got nil")
		}
	})
}
