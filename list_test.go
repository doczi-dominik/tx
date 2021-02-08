package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLoadLocal(t *testing.T) {
	// Prepare data to be loaded
	InitTestingPathVariables(t)

	taskfile := CreateTaskfile(TaskfilePath)
	taskfile.WriteString("a\nb\nc\n")
	taskfile.Close()

	// Simulate a tasklist
	tl := &Tasklist{
		filePath: TaskfilePath,
		tasks:    make(map[int]Task),
	}

	tl.LoadLocal()

	AssertEqual(t, tl.loaded, true, "Simulated tasklist's loaded field is false")
	AssertEqual(t, len(tl.tasks), 3, "Simulated tasklist does not have 3 tasks in it")
}

func TestSaveLocal(t *testing.T) {
	InitTestingPathVariables(t)

	ConfigOptions.DeleteIfEmpty = true

	now := time.Now()

	// Simulate a tasklist with data in it
	tl := &Tasklist{
		filePath: TaskfilePath,
		tasks: map[int]Task{
			1: {text: "a", hash: "idc", creationDate: now, finishedDate: time.Unix(0, 0)},
		},
	}

	// Check serialization
	tl.SerializeTasks()

	AssertNotEqual(t, len(tl.serialized), 0, "Simulated tasklist's serialized field is empty")

	// Save then check file contents
	tl.SaveLocal()

	var contents []byte

	scanner := bufio.NewScanner(OpenTaskfile(TaskfilePath, false))

	for scanner.Scan() {
		contents = append(contents, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		Error(ErrTaskfileRead, TaskfilePath, err)
	}

	checkContents := fmt.Sprintf("a | id:idc, creation:%s, finished:1970/01/01/01/00", now.Format(FullDateFormat))

	AssertEqual(t, string(contents), checkContents, "Taskfile's contents do not match the manually constructed one")

	// Test --delete-if-empty
	tl.tasks = make(map[int]Task)
	tl.SerializeTasks()

	// Save and test if the file still exists
	tl.SaveLocal()

	if _, err := os.Stat(TaskfilePath); err == nil {
		t.Fatal("Empty taskfile was not deleted")
	}

	// Test if backup was created
	if _, err := os.Stat(BackupfilePath); os.IsNotExist(err) {
		t.Fatal("Backup file was not created after SaveLocal()")
	}
}
