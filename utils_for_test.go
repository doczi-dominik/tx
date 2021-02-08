package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

// InitTestingEnv initializes a tasklist with a set of 7 default "hello world"
// tasks.
func InitTestingEnv(tl **Tasklist) {
	*tl = &Tasklist{
		loaded: true,
		tasks: map[int]Task{
			1: {text: "hello world"},
			2: {text: "hello world"},
			3: {text: "hello world"},
			4: {text: "hello world"},
			5: {text: "hello world"},
			6: {text: "hello world"},
			7: {text: "hello world"},
		},
	}
}

// InitNumberedTestingEnv initializes a tasklist with 7 tasks corresponding to
// the lowercase english words for 1 to 7.
func InitNumberedTestingEnv(tl **Tasklist) {
	*tl = &Tasklist{
		loaded: true,
		tasks: map[int]Task{
			1: {text: "one"},
			2: {text: "two"},
			3: {text: "three"},
			4: {text: "four"},
			5: {text: "five"},
			6: {text: "six"},
			7: {text: "seven"},
		},
	}
}

// InitEmptyTestingEnv initializes a tasklist with an empty tasks map.
func InitEmptyTestingEnv(tl **Tasklist) {
	*tl = &Tasklist{
		loaded: true,
		tasks:  map[int]Task{},
	}
}

// InitTestingPathVariables initializes all filepaths to a temporary directory
// managed by the testing library.
func InitTestingPathVariables(t *testing.T) {
	// Circumventing the no re-initialization rule for testing.
	SyncfilePath = ""
	InitPathVariables(t.TempDir() + "/tasks")
}

// AssertEqual fails the current test if two values are not equal.
func AssertEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	if a != b {
		t.Fatalf("%s\n%v != %v", msg, a, b)
	}
}

// AssertNotEqual fails the current test if two values are equal.
func AssertNotEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	if a == b {
		t.Fatalf("%s\n%v == %v", msg, a, b)
	}
}

// AssertTaskText fails the current test if a task's text field does not
// match the provided one.
func AssertTaskText(t *testing.T, task Task, text string) {
	msg := fmt.Sprintf("Task's text does not equal \"%s\"", text)
	AssertEqual(t, task.text, text, msg)
}

// AssertTaskText fails the current test if a task's hash field does not
// match the provided one.
func AssertTaskHash(t *testing.T, task Task, hash string) {
	AssertEqual(t, task.hash, hash, "Task's hash field does not equal its text field's SHA-1 hash")
}

// AssertTaskText fails the current test if a task's creationDate field does
// not match the provided one.
func AssertTaskCreationDate(t *testing.T, task Task, date time.Time) {
	AssertEqual(t, task.creationDate, date, "Task's creationDate field does not equal the provided date")
}

// AssertTaskText fails the current test if a task's finishedDate field does
// not match the provided one.
func AssertTaskFinishedDate(t *testing.T, task Task, date time.Time) {
	AssertEqual(t, task.finishedDate, date, "Task's finishedDate field does not equal the provided date")
}

// AssertMainTaskText fails the current test if the selected task's text does
// not match the provided text.
func AssertMainTaskText(t *testing.T, index int, text string) {
	msg := fmt.Sprintf("Text of active task at index %d does not match the provided text", index)
	AssertEqual(t, MainList.tasks[index].text, text, msg)
}

// AssertDoneTaskText fails the current test if the selected task's text does
// not match the provided text.
func AssertDoneTaskText(t *testing.T, index int, text string) {
	msg := fmt.Sprintf("Text of finished task at index %d does not match the provided text", index)
	AssertEqual(t, DoneList.tasks[index].text, text, msg)
}

// AssertDeletedTask fails the current test if a task is available at
// the provided index.
func AssertDeletedTask(t *testing.T, index int, tl *Tasklist) {
	_, ok := tl.tasks[index]

	if ok {
		t.Fatalf("Task accessible at index %d", index)
	}
}

// AssertEmptyTasklist fails the current test if the provided tasklist is not
// empty.
func AssertEmptyTasklist(t *testing.T, tl *Tasklist) {
	if !tl.IsEmpty() {
		var name string

		if tl == MainList {
			name = "MainList"
		} else {
			name = "DoneList"
		}

		t.Fatalf("%s should be empty!", name)
	}
}

// AssertNonEmptyTasklist fails the current test if the provided tasklist is
// empty.
func AssertNonEmptyTasklist(t *testing.T, tl *Tasklist) {
	if tl.IsEmpty() {
		var name string

		if tl == MainList {
			name = "MainList"
		} else {
			name = "DoneList"
		}

		t.Fatalf("%s should not be empty!", name)
	}
}

// AssertExitError fails the current test if a test exits with an exit code
// not equal to the provided one.
func AssertExitError(t *testing.T, testName string, wantedExitCode int, test func()) {
	if os.Getenv("TX_TESTING") == "true" {
		test()
		return
	}

	commandLine := []string{os.Args[0], "-test.run=" + testName}

	cmd := exec.Command(commandLine[0], commandLine[1:]...)
	cmd.Env = append(os.Environ(), "TX_TESTING=true")

	msg := "Child process exited with code %d (wanted code %d)"

	output, err := cmd.CombinedOutput()

	if err == nil {
		if wantedExitCode != -1 {
			t.Fatal(msg, 0, wantedExitCode)
		}
		return
	}

	exitError := err.(*exec.ExitError)
	exitCode := exitError.ExitCode()

	if wantedExitCode+1 != exitCode {
		t.Fatalf(msg+"\nSTDOUT+STDERR: %s\n", exitCode, wantedExitCode+1, output)
	}
}

// AssertExitSuccess fails the current test if a test exits with an exit
// code other than 0.
func AssertExitSuccess(t *testing.T, testName string, test func()) {
	AssertExitError(t, testName, -1, test)
}
