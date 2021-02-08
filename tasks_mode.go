package main

import (
	"regexp"
	"strings"
	"time"
)

// TaskActions contains all actions and positional arguments for task
// management mode.
type TaskActions struct {
	Add      func(string) `short:"a" long:"add" description:"Add a new task. Use when specifying multiple actions." value-name:"TEXT"`
	Edit     func(string) `short:"e" long:"edit" description:"Replace an entire task/words from a task using sed syntax" value-name:"<SELECT,TEXT or SELECT/OLD/NEW>"`
	Finish   func(string) `short:"f" long:"finish" description:"Mark TASK as finished" value-name:"SELECT"`
	Remove   func(string) `short:"r" long:"remove" description:"Remove TASK from list" value-name:"SELECT"`
	Wipe     func()       `short:"w" long:"wipe" description:"Remove all tasks"`
	Complete func()       `short:"c" long:"complete" description:"Mark all tasks as finished"`

	Args struct {
		TEXT []string
	} `positional-args:"yes"`
}

var taskActions TaskActions

// Execute runs after all actions; adds positional arguments as tasks, lists
// tasks and saves the tasklist if necessary.
func (a *TaskActions) Execute(args []string) error {
	ListManager.EnsureInitialized(MainList)

	// Join positional arguments and add them as a
	// new task.
	if len(a.Args.TEXT) != 0 {
		text := strings.Join(a.Args.TEXT, " ")
		MainList.Add(NewTask(text))
	}

	MainList.Show(OutputOptions.Format)

	ListManager.Save()
	return nil
}

// add adds a new task to the task list.
func add(text string) {
	ListManager.EnsureInitialized(MainList)

	MainList.Add(NewTask(text))
}

// edit recognizes two formats (full replace and sed-style) and edits the
// provided task accordingly.
func edit(cmd string) {
	ListManager.EnsureInitialized(MainList)
	exitOnEmptyTasks("edit")

	var parts []string
	recursiveSplit(cmd, &parts)

	l := len(parts)
	if 2 > l || l > 3 {
		Error(ErrEditInvalidSelector, cmd)
	}

	index := MainList.InterpretSelectorPart(parts[0], MainList.OrderKeys())
	oldTask, exists := MainList.tasks[index]

	if !exists {
		Error(ErrInvalidIndex, "Edit", index)
	}

	var newText string

	if l == 2 {
		newText = strings.ReplaceAll(parts[1], `\/`, `/`)
	} else {
		search := strings.ReplaceAll(parts[1], `\/`, `/`)
		repl := strings.ReplaceAll(parts[2], `\/`, `/`)

		newText = strings.ReplaceAll(oldTask.text, search, repl)
	}

	newTask := Task{
		text:         newText,
		creationDate: oldTask.creationDate,
		finishedDate: oldTask.finishedDate,
	}

	err := newTask.Validate()

	if err != nil {
		Error(ErrTaskValidation, "Edit", err)
	}

	MainList.tasks[index] = Task{
		text:         newText,
		creationDate: oldTask.creationDate,
		finishedDate: oldTask.finishedDate,
	}

	MainList.MarkModified()
}

// finish initializes a "finished tasks" list, adds tasks to it, then removes
// the tasks and writes the donelist to file.
func finish(selector string) {
	ListManager.EnsureInitialized(MainList)
	exitOnEmptyTasks("Finish")
	ListManager.EnsureInitialized(DoneList)

	indexes, err := MainList.SelectTasks(selector)

	if err != nil {
		Error(ErrInvalidSelector, "Finish", selector)
	}

	for _, i := range indexes {
		task := MainList.tasks[i]

		task.finishedDate = time.Now()
		DoneList.Add(task)
	}

	MainList.Remove(indexes)
}

// remove removes a task from the tasklist.
func remove(selector string) {
	ListManager.EnsureInitialized(MainList)
	exitOnEmptyTasks("Remove")

	indexes, err := MainList.SelectTasks(selector)

	if err != nil {
		Error(ErrInvalidSelector, "Remove", selector)
	}

	MainList.Remove(indexes)
}

// wipeTasks is a convenience action that removes every task from the tasklist.
func wipeTasks() {
	ListManager.EnsureInitialized(MainList)
	exitOnEmptyTasks("Wipe")

	MainList.tasks = make(map[int]Task)

	MainList.MarkModified()
}

// complete is a convenience action that marks all tasks as finished.
func complete() {
	ListManager.EnsureInitialized(MainList)
	exitOnEmptyTasks("Complete")

	finish("f-l")
}

// init gets called when the package is imported; assigns functions to the
// respective action and adds the subcommand to the global argument parser.
func init() {
	taskActions.Add = add
	taskActions.Edit = edit
	taskActions.Finish = finish
	taskActions.Remove = remove
	taskActions.Wipe = wipeTasks
	taskActions.Complete = complete

	GlobalParser.AddCommand("tasks", "Manage active tasks", "", &taskActions)
}

func exitOnEmptyTasks(caller string) {
	if MainList.IsEmpty() {
		Error(ErrTasklistEmpty, caller)
	}
}

// Shoutout to RE2 for not including (optional) lookbehind support
func recursiveSplit(source string, parts *[]string) {
	src := strings.TrimSpace(source)

	if src == "" {
		return
	}

	firstSeparatorIndex := regexp.MustCompile(`[^\\]\/`).FindStringIndex(src)

	if len(firstSeparatorIndex) == 0 {
		*parts = append(*parts, src)
		return
	}

	// [0] + 1 used since pattern consumes one character to check if it is
	// escaped or not.
	start := firstSeparatorIndex[0] + 1

	*parts = append(*parts, src[:start])

	recursiveSplit(src[start+1:], parts)
}
