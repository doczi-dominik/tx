package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// ConfigOptions holds all the global configuration options that apply to both
// active tasklist management and finished tasklist management.
var ConfigOptions struct {
	List            string `short:"L" long:"list" description:"Path to the regular taskfile" value-name:"PATH"`
	DeleteIfEmpty   bool   `short:"D" long:"delete-if-empty" description:"Delete the taskfile if it becomes empty"`
	Callback        string `short:"C" long:"callback" description:"Path to script/command (+ args) to run after writing a tasklist" value-name:"CMD"`
	Offline         bool   `short:"O" long:"offline" description:"Disable loading from network and default to local tasklists only"`
	Reckless        bool   `short:"R" long:"reckless" description:"Disable taking local backups after modifying a taskfile"`
	Quiet           bool   `short:"Q" long:"quiet" description:"Disable the printing of warning messages"`
	FallbackSyncURL string `short:"U" long:"fallback-sync-url" description:"The URL of the Sync service to use if no explicit URL is specified for the tasklist." value-name:"URL"`
}

// OutputOptions holds all the options which modify the output.
var OutputOptions struct {
	Format string `short:"o" long:"format" description:"Defines the output format.\nPlaceholders: {index}, {task}, {creationTime}, {creationDate}, {finishedTime}, {finishedDate}" value-name:"STRING"`
}

// RunCallback executes the configured callback command (if any).
func RunCallback() {
	callback := strings.TrimSpace(ConfigOptions.Callback)

	if callback == "" {
		return
	}

	cmdLine := strings.Split(callback, " ")

	command := exec.Command(cmdLine[0], cmdLine[1:]...)
	output, err := command.CombinedOutput()

	if err != nil {
		Error(ErrCallback, callback, err)
	}

	if len(output) != 0 {
		fmt.Printf("\n%s", string(output))
	}
}

func init() {
	// Set defaults manually so they are available before parsing finishes.
	ConfigOptions.List = "tasks"
	ConfigOptions.FallbackSyncURL = ""
	OutputOptions.Format = "{index} - {task}"

	GlobalParser.AddGroup("Configuration Options", "", &ConfigOptions)
	GlobalParser.AddGroup("Output Options", "", &OutputOptions)
}
