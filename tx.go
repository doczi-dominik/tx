package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

// GlobalParser is the argument parser which holds all options and arguments.
var GlobalParser *flags.Parser = flags.NewNamedParser("tx", flags.HelpFlag|flags.PassDoubleDash)
var GlobalFS FS

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

func main() {
	GlobalFS = createFileFS()

	_, err := GlobalParser.Parse()

	if err != nil {
		if flags.WroteHelp(err) {
			os.Stderr.WriteString(err.Error())
		} else {
			Error(ErrFlagParsing, err)
		}
	}
}
