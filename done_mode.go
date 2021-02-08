package main

// DoneActions contains all actions for finished task management.
type DoneActions struct {
	Restore    func(string) `short:"r" long:"restore" description:"Restore a task to the main tasklist" value-name:"SELECT"`
	RestoreAll func()       `short:"a" long:"restore-all" description:"Restores all finished tasks"`
	Delete     func(string) `short:"d" long:"delete" description:"Remove a finished task from the list" value-name:"SELECT"`
	Wipe       func()       `short:"w" long:"wipe" description:"Remove all tasks"`
}

var doneActions DoneActions

// Execute runs after all actions; lists finished tasks and writes them to file
// if necessary.
func (a *DoneActions) Execute(args []string) error {
	ListManager.EnsureInitialized(DoneList)

	DoneList.Show(OutputOptions.Format)

	ListManager.Save()

	return nil
}

func restore(selector string) {
	ListManager.EnsureInitialized(DoneList)
	exitOnEmptyDone("Restore")
	ListManager.EnsureInitialized(MainList)

	indexes, err := DoneList.SelectTasks(selector)

	if err != nil {
		Error(ErrInvalidSelector, "Restore", selector)
	}

	for _, i := range indexes {
		MainList.Add(DoneList.tasks[i])
	}

	DoneList.Remove(indexes)
}

func restoreAll() {
	ListManager.EnsureInitialized(DoneList)
	exitOnEmptyDone("Restore All")

	restore("f-l")
}

func deleteDone(selector string) {
	ListManager.EnsureInitialized(DoneList)
	exitOnEmptyDone("Delete")

	indexes, err := DoneList.SelectTasks(selector)

	if err != nil {
		Error(ErrInvalidSelector, "Delete", selector)
	}

	DoneList.Remove(indexes)
}

func wipeDone() {
	ListManager.EnsureInitialized(DoneList)
	exitOnEmptyDone("Wipe")

	DoneList.tasks = make(map[int]Task)

	DoneList.MarkModified()
}

// init gets called when the package is imported; assigns functions to the
// respective action and adds the subcommand to the global argument parser.
func init() {
	doneActions.Restore = restore
	doneActions.RestoreAll = restoreAll
	doneActions.Delete = deleteDone
	doneActions.Wipe = wipeDone

	GlobalParser.AddCommand("done", "Manage finished tasks", "", &doneActions)
}

func exitOnEmptyDone(caller string) {
	if DoneList.IsEmpty() {
		Error(ErrTasklistEmpty, caller)
	}
}
