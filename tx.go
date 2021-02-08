package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

// GlobalParser is the argument parser which holds all options and arguments.
var GlobalParser *flags.Parser = flags.NewNamedParser("tx", flags.HelpFlag|flags.PassDoubleDash)

var (
	// TaskfilePath holds the path to the current active taskfile. By default,
	// it's "./tasks".
	TaskfilePath string
	// DonefilePath holds the path to the current finished taskfile. The path
	// is derived from TaskfilePath like so: "./.{TaskfilePath}.done".
	DonefilePath string
	// SyncfilePath holds the path to the current Syncfile. The path is derived
	// from TaskfilePath like so: "./.{TaskfilePath}.sync".
	SyncfilePath string
	// BackupfilePath holds the path to the current backup taskfile. The path
	// is derived from TaskfilePath like so: "./.{TaskfilePath}.bak".
	BackupfilePath string
)

var (
	// MainList is the tasklist which holds all active tasks. Managed by
	// ListManager.
	MainList = &Tasklist{}
	// DoneList is the tasklist which holds all finished tasks. Managed by
	// ListManager.
	DoneList = &Tasklist{}
	// ListManager is responsible for loading and saving taskfiles, both local
	// files and Sync service data.
	ListManager = &TasklistManager{}
)

// InitPathVariables initializes all filepath variables. Re-initialization is
// considered irregular and the function will return.
func InitPathVariables(taskfilePath string) {
	if SyncfilePath != "" {
		return
	}

	TaskfilePath = taskfilePath
	DonefilePath = GetMetafilePath(".done")
	SyncfilePath = GetMetafilePath(".sync")
	BackupfilePath = GetMetafilePath(".bak")
}

func main() {
	_, err := GlobalParser.Parse()

	if err != nil {
		if flags.WroteHelp(err) {
			os.Stderr.WriteString(err.Error())
		} else {
			Error(ErrFlagParsing, err)
		}
	}
}
