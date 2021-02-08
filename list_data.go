package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// MarkModified sets the underlying modfied flag to true.
func (tl *Tasklist) MarkModified() {
	tl.modified = true
}

// IsEmpty returns whether the tasklist is empty.
func (tl *Tasklist) IsEmpty() (empty bool) {
	return len(tl.tasks) == 0
}

// OrderKeys sorts the keys of the tasklist for enumarating.
func (tl *Tasklist) OrderKeys() (keys []int) {
	for key := range tl.tasks {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	return
}

// Add adds a task to the tasklist.
func (tl *Tasklist) Add(newTask Task) {
	err := newTask.Validate()

	if err != nil {
		Error(ErrTaskValidation, "Add", err)
	}

	// Do not calculate index if tasklist is empty
	if tl.IsEmpty() {
		tl.tasks[1] = newTask
	} else {
		indexes := tl.OrderKeys()
		lastIndex := indexes[len(indexes)-1]

		tl.tasks[lastIndex+1] = newTask
	}

	// Do not mark as modified when loading tasks, only when
	// adding new ones.
	if tl.loaded {
		tl.MarkModified()
	}
}

// Remove removes one or more tasks from the tasklist.
func (tl *Tasklist) Remove(indexes []int) {
	for _, index := range indexes {
		delete(tl.tasks, index)
	}

	if len(indexes) > 0 {
		tl.MarkModified()
	}
}

// Show generates and prints the output according to a user-provided format
// string.
func (tl *Tasklist) Show(format string) {
	if tl.IsEmpty() {
		return
	}

	keys := tl.OrderKeys()

	padding := fmt.Sprintf("%%%dd", 1+len(keys)/10)

	for displayIndex, index := range keys {
		task := tl.tasks[index]

		creationDate := task.creationDate.Format(DateFormat)
		creationTime := task.creationDate.Format(DisplayTimeFormat)

		finishedDate := task.finishedDate.Format(DateFormat)
		finishedTime := task.finishedDate.Format(DisplayTimeFormat)

		replacer := strings.NewReplacer(
			"{index}", fmt.Sprintf(padding, displayIndex+1),
			"{creationDate}", creationDate,
			"{creationTime}", creationTime,
			"{finishedDate}", finishedDate,
			"{finishedTime}", finishedTime,
			"{task}", task.text,
			"{{", "{",
			"}}", "}",
		)

		fmt.Println(replacer.Replace(format))
	}
}

// InterpretSelectorPart converts a selector token to a concrete index.
func (tl *Tasklist) InterpretSelectorPart(part string, keys []int) (result int) {
	s := strings.ToLower(part)

	if s == "f" || s == "r" {
		return keys[0]
	}

	if s == "l" {
		return keys[len(keys)-1]
	}

	result, _ = strconv.Atoi(s)
	return
}

// SelectTasks converts a "selector" string to a slice of valid indexes.
func (tl *Tasklist) SelectTasks(selector string) (indexes []int, err error) {
	selector = strings.TrimSpace(selector)
	keys := tl.OrderKeys()

	// Range notation
	indexRange := regexp.MustCompile(`([rfl]|\d+)-([rfl]|\d+)`).FindStringSubmatch(selector)

	if len(indexRange) == 3 {

		valueOne := tl.InterpretSelectorPart(indexRange[1], keys)
		valueTwo := tl.InterpretSelectorPart(indexRange[2], keys)

		start := valueOne
		stop := valueTwo

		// Reverse range if needed
		if valueOne > valueTwo {
			start = valueTwo
			stop = valueOne
		}

		for i := start; i <= stop; i++ {
			_, ok := tl.tasks[i]

			if ok {
				indexes = append(indexes, i)
			}
		}

		return
	}

	// CSV / Regular notation
	indexList := regexp.MustCompile(`(?i)(?:(?:[rfl]|\d+)[ \t]*,*[ \t]*)+`)

	if indexList.MatchString(selector) {
		stringList := regexp.MustCompile(`,`).Split(selector, -1)

		for _, str := range stringList {
			indexes = append(indexes, tl.InterpretSelectorPart(str, keys))
		}

		return
	}

	return indexes, fmt.Errorf("invalidFormat")
}

// ParseTasklines tries to parse tasks and add them to the tasklists.
// If configured, incorrectly formatted tasks are collected and displayed.
func (tl *Tasklist) ParseTasklines(source string, reader io.Reader) {
	var (
		errorOccured bool
		lineNumber   int
	)

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		newTask, err := ParseTask(line)

		if err != nil {
			if err.Error() == "ignore" {
				continue
			}

			if err.Error() == "writeNewMeta" {
				tl.MarkModified()
			}
		}

		tl.Add(newTask)
	}

	if err := scanner.Err(); err != nil {
		Error(ErrTaskfileRead, source, err)
	}

	// Separate corrupt tasklines from regular ones.
	if errorOccured {
		fmt.Println()
	}
}

// SerializeTasks serializes all tasks and appends them to the `serialized`
// field of the tasklist.
func (tl *Tasklist) SerializeTasks() {
	for _, index := range tl.OrderKeys() {
		task := tl.tasks[index]
		tl.serialized = append(tl.serialized, task.Serialize()...)
	}
}
