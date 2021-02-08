package main

import "testing"

func TestRestore(t *testing.T) {
	InitNumberedTestingEnv(&DoneList)
	InitEmptyTestingEnv(&MainList)

	t.Run("index", func(t *testing.T) {
		restore("3")

		AssertDeletedTask(t, 3, DoneList)

		AssertMainTaskText(t, 1, "three")
	})

	t.Run("letter", func(t *testing.T) {
		restore("l")

		AssertDeletedTask(t, 7, DoneList)

		AssertMainTaskText(t, 2, "seven")
	})

	t.Run("csv", func(t *testing.T) {
		restore("2,4")

		AssertDeletedTask(t, 2, DoneList)
		AssertDeletedTask(t, 4, DoneList)

		AssertMainTaskText(t, 3, "two")
		AssertMainTaskText(t, 4, "four")
	})

	t.Run("range", func(t *testing.T) {
		restore("1-6")

		AssertDeletedTask(t, 1, DoneList)
		AssertDeletedTask(t, 5, DoneList)
		AssertDeletedTask(t, 6, DoneList)

		AssertMainTaskText(t, 5, "one")
		AssertMainTaskText(t, 6, "five")
		AssertMainTaskText(t, 7, "six")
	})

	AssertEmptyTasklist(t, DoneList)
	AssertEqual(t, len(MainList.tasks), 7, "MainList does not have 7 tasks")
}

func TestRestoreAll(t *testing.T) {
	InitTestingEnv(&DoneList)
	InitEmptyTestingEnv(&MainList)

	restoreAll()

	AssertEmptyTasklist(t, DoneList)
	AssertEqual(t, len(MainList.tasks), 7, "MainList does not have 7 tasks")
}

func TestDelete(t *testing.T) {
	InitTestingEnv(&DoneList)

	t.Run("index", func(t *testing.T) {
		deleteDone("3")
		AssertDeletedTask(t, 3, DoneList)
	})

	t.Run("letter", func(t *testing.T) {
		deleteDone("l")
		AssertDeletedTask(t, 7, DoneList)
	})

	t.Run("csv", func(t *testing.T) {
		deleteDone("2,4")
		AssertDeletedTask(t, 2, DoneList)
		AssertDeletedTask(t, 4, DoneList)
	})

	t.Run("range", func(t *testing.T) {
		deleteDone("1-6")
		AssertDeletedTask(t, 1, DoneList)
		AssertDeletedTask(t, 5, DoneList)
		AssertDeletedTask(t, 6, DoneList)
	})

	AssertEmptyTasklist(t, DoneList)
}
