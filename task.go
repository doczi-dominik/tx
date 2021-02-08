package main

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

// Task represents text as content a creation date and a date indicating when
// it was marked as complete.
type Task struct {
	text         string
	creationDate time.Time
	finishedDate time.Time
	hash         string
}

// Validate checks if a task contains valid data and corrects common errors.
func (t *Task) Validate() error {
	t.text = strings.TrimSpace(t.text)

	if strings.Contains(t.text, "\n") {
		return fmt.Errorf("Task cannot contain newline")
	}

	return nil
}

// Serialize converts a Task object to a taskline that can be written to a
// taskfile.
func (t *Task) Serialize() (data []byte) {
	escapedText := strings.TrimSpace(strings.ReplaceAll(t.text, `|`, `\|`))

	creationString := t.creationDate.Format(FullDateFormat)
	finishedString := t.finishedDate.Format(FullDateFormat)

	line := fmt.Sprintf("%s | id:%s, creation:%s, finished:%s\n", escapedText, t.hash, creationString, finishedString)

	data = []byte(line)
	return
}

// ParseTask takes a line from a taskfile and creates a Task object from it.
func ParseTask(line string) (newTask Task, err error) {
	line = strings.TrimSpace(line)
	newTask.creationDate = time.Now()
	newTask.finishedDate = time.Unix(0, 0)

	if line == "" {
		err = fmt.Errorf("ignore")
		return
	}

	// If the line does not contain a part separator, interpret it as a new
	// task.
	if !SeparatorPattern.MatchString(line) {
		newTask = NewTask(strings.ReplaceAll(line, `\|`, `|`))

		err = fmt.Errorf("writeNewMeta")

		return
	}

	// Split task text and syncdata on separator '|'
	// ----------------------------
	// [0] + 1 used since the pattern consumes the preceding character in order
	// to check if it is escaped
	separatorLocation := SeparatorPattern.FindStringIndex(line)[0] + 1

	newTask.text = strings.TrimSpace(strings.ReplaceAll(line[:separatorLocation], `\|`, `|`))
	metadata := strings.TrimSpace(line[separatorLocation+1:])

	parsedCrDate := CreationDatePattern.FindStringSubmatch(metadata)

	if len(parsedCrDate) != 0 {
		creationDate, err := time.Parse(FullDateFormat, parsedCrDate[1])

		if err == nil {
			newTask.creationDate = creationDate
		} else {
			Warn("Could not parse creation date: \"%s\". Using current time.", parsedCrDate[1])
		}
	} else {
		err = fmt.Errorf("writeNewMeta")
	}

	parsedFiDate := FinishedDatePattern.FindStringSubmatch(metadata)

	if len(parsedFiDate) != 0 {
		finishedDate, err := time.Parse(FullDateFormat, parsedFiDate[1])

		if err == nil {
			newTask.finishedDate = finishedDate
		} else {
			Warn("Could not parse finished date: \"%s\". Using Unix epoch.")
		}
	} else {
		err = fmt.Errorf("writeNewMeta")
	}

	parsedHash := HashPattern.FindStringSubmatch(metadata)

	if len(parsedHash) != 0 {
		newTask.hash = parsedHash[1]
	} else {
		newTask.hash = hexHash(newTask.text)
		err = fmt.Errorf("writeNewMeta")
	}

	return
}

// NewTask creates a task and sets default values.
func NewTask(text string) (newTask Task) {
	newTask.text = text
	newTask.creationDate = time.Now()
	newTask.finishedDate = time.Unix(0, 0)
	newTask.hash = hexHash(text)

	return
}

func hexHash(s string) (hash string) {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}
