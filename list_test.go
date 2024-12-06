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
	taskfile := GlobalFS.createTaskfile()
	taskfile.WriteString("a\nb\nc\n")
	taskfile.Close()

	// Simulate a tasklist
	tl := &Tasklist{
		filePath: GlobalFS.getTaskfilePath(),
		tasks:    make(map[int]Task),
	}

	tl.LoadLocal()

	AssertEqual(t, tl.loaded, true, "Simulated tasklist's loaded field is false")
	AssertEqual(t, len(tl.tasks), 3, "Simulated tasklist does not have 3 tasks in it")
}

func TestSaveLocal(t *testing.T) {
	doTest := func() {
		now := time.Now()
		unixBirth := time.Unix(0, 0)

		tfPath := GlobalFS.getTaskfilePath()

		// Simulate a tasklist with data in it
		tl := &Tasklist{
			filePath: tfPath,
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

		scanner := bufio.NewScanner(GlobalFS.openTaskfile(false))

		for scanner.Scan() {
			contents = append(contents, scanner.Bytes()...)
		}

		if err := scanner.Err(); err != nil {
			Error(ErrTaskfileRead, tfPath, err)
		}

		checkContents := fmt.Sprintf("a | id:idc, creation:%s, finished:%s", now.Format(FullDateFormat), unixBirth.Format(FullDateFormat))

		AssertEqual(t, string(contents), checkContents, "Taskfile's contents do not match the manually constructed one")

		// Test --delete-if-empty
		tl.tasks = make(map[int]Task)
		tl.SerializeTasks()

		// Save and test if the file still exists
		tl.SaveLocal()

		// Test if backup was created
		if _, err := GlobalFS.Stat(GetMetafilePath(".bak", tfPath)); os.IsNotExist(err) {
			t.Fatal("Backup file was not created after SaveLocal()")
		}
	}

	t.Run("deleteIfEmpty-true", func(t *testing.T) {
		ConfigOptions.DeleteIfEmpty = true

		doTest()

		tfPath := GlobalFS.getTaskfilePath()
		_, err := GlobalFS.Stat(tfPath)

		if err == nil {
			t.Fatal("Empty taskfile was not deleted")
		}
	})

	t.Run("deleteIfEmpty-false", func(t *testing.T) {
		ConfigOptions.DeleteIfEmpty = false

		doTest()
	})
}
