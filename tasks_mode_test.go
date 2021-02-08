package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	InitEmptyTestingEnv(&MainList)

	t.Run("regular", func(t *testing.T) {
		add("a new task")
		AssertMainTaskText(t, 1, "a new task")
	})

	t.Run("invalid", func(t *testing.T) {
		AssertExitError(t, "TestAdd/invalid", ErrTaskValidation, func() {
			add("this\nshould\nbe\ninvalid")
		})
	})
}

func TestEdit(t *testing.T) {
	InitTestingEnv(&MainList)

	t.Run("full", func(t *testing.T) {
		edit("1/changed")
		AssertMainTaskText(t, 1, "changed")
	})

	t.Run("sed", func(t *testing.T) {
		edit("2/world/people")
		AssertMainTaskText(t, 2, "hello people")
	})

	t.Run("full_trailing", func(t *testing.T) {
		edit("3/logic")
		AssertMainTaskText(t, 3, "logic")
	})

	t.Run("sed_trailing", func(t *testing.T) {
		edit("4/hello/greetings/")
		AssertMainTaskText(t, 4, "greetings world")
	})

	t.Run("full_escaped", func(t *testing.T) {
		edit(`5/hello\/world`)
		AssertMainTaskText(t, 5, "hello/world")
	})

	t.Run("sed_escaped", func(t *testing.T) {
		edit(`6/world/world\/worlds/`)
		AssertMainTaskText(t, 6, "hello world/worlds")
	})

	t.Run("letter_selector", func(t *testing.T) {
		edit("l/seven")
		AssertMainTaskText(t, 7, "seven")
	})

	t.Run("invalid_index", func(t *testing.T) {
		AssertExitError(t, "TestEdit/invalid_index", ErrInvalidIndex, func() {
			edit("9/non-existent")
		})
	})
}

func TestFinish(t *testing.T) {
	InitNumberedTestingEnv(&MainList)
	InitEmptyTestingEnv(&DoneList)

	t.Run("index", func(t *testing.T) {
		finish("3")

		AssertDeletedTask(t, 3, MainList)

		AssertDoneTaskText(t, 1, "three")
	})

	t.Run("letter", func(t *testing.T) {
		finish("l")

		AssertDeletedTask(t, 7, MainList)

		AssertDoneTaskText(t, 2, "seven")
	})

	t.Run("csv", func(t *testing.T) {
		finish("2,4")

		AssertDeletedTask(t, 2, MainList)
		AssertDeletedTask(t, 4, MainList)

		AssertDoneTaskText(t, 3, "two")
		AssertDoneTaskText(t, 4, "four")
	})

	t.Run("range", func(t *testing.T) {
		finish("1-6")

		AssertDeletedTask(t, 1, MainList)
		AssertDeletedTask(t, 5, MainList)
		AssertDeletedTask(t, 6, MainList)

		AssertDoneTaskText(t, 5, "one")
		AssertDoneTaskText(t, 6, "five")
		AssertDoneTaskText(t, 7, "six")
	})
}

func TestRemove(t *testing.T) {
	InitTestingEnv(&MainList)

	t.Run("index", func(t *testing.T) {
		remove("3")
		AssertDeletedTask(t, 3, MainList)
	})

	t.Run("letter", func(t *testing.T) {
		remove("l")
		AssertDeletedTask(t, 7, MainList)
	})

	t.Run("csv", func(t *testing.T) {
		remove("2,4")
		AssertDeletedTask(t, 2, MainList)
		AssertDeletedTask(t, 4, MainList)
	})

	t.Run("range", func(t *testing.T) {
		remove("1-6")
		AssertDeletedTask(t, 1, MainList)
		AssertDeletedTask(t, 5, MainList)
		AssertDeletedTask(t, 6, MainList)
	})

	AssertEmptyTasklist(t, MainList)
}

func TestWipe(t *testing.T) {
	InitTestingEnv(&MainList)

	wipeTasks()

	AssertEmptyTasklist(t, MainList)
}

func TestComplete(t *testing.T) {
	InitTestingEnv(&MainList)
	InitEmptyTestingEnv(&DoneList)

	complete()

	AssertEmptyTasklist(t, MainList)
	AssertNonEmptyTasklist(t, DoneList)
}
