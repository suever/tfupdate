package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/minamijoyo/tfupdate/command"
	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
)

// Version is a version number.
var version = "0.2.0"

// UI is a user interface which is a global variable for mocking.
var UI cli.Ui

func init() {
	UI = &cli.BasicUi{
		Writer: os.Stdout,
	}
}

func main() {
	log.SetOutput(logOutput())
	log.Printf("[INFO] CLI args: %#v", os.Args)

	commands := initCommands()

	args := os.Args[1:]

	c := &cli.CLI{
		Name:                  "tfupdate",
		Version:               version,
		Args:                  args,
		Commands:              commands,
		HelpWriter:            os.Stdout,
		Autocomplete:          true,
		AutocompleteInstall:   "install-autocomplete",
		AutocompleteUninstall: "uninstall-autocomplete",
	}

	exitStatus, err := c.Run()
	if err != nil {
		UI.Error(fmt.Sprintf("Failed to execute CLI: %s", err))
	}

	os.Exit(exitStatus)
}

func logOutput() io.Writer {
	levels := []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}
	minLevel := os.Getenv("TFUPDATE_LOG")

	// default log writer is null device.
	writer := ioutil.Discard
	if minLevel != "" {
		writer = os.Stderr
	}

	filter := &logutils.LevelFilter{
		Levels:   levels,
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   writer,
	}

	return filter
}

func initCommands() map[string]cli.CommandFactory {
	meta := command.Meta{
		UI: UI,
		Fs: afero.NewOsFs(),
	}

	commands := map[string]cli.CommandFactory{
		"terraform": func() (cli.Command, error) {
			return &command.TerraformCommand{
				Meta: meta,
			}, nil
		},
		"provider": func() (cli.Command, error) {
			return &command.ProviderCommand{
				Meta: meta,
			}, nil
		},
		"release": func() (cli.Command, error) {
			return &command.ReleaseCommand{
				Meta: meta,
			}, nil
		},
		"release latest": func() (cli.Command, error) {
			return &command.ReleaseLatestCommand{
				Meta: meta,
			}, nil
		},
	}

	return commands
}
