package main

import (
	"testing"
	"time"
)

func TestParseTask(t *testing.T) {
	t.Run("empty_line", func(t *testing.T) {
		_, err := ParseTask("")

		AssertEqual(t, err.Error(), "ignore", "Empty line is not ignored")
	})

	t.Run("no_separator", func(t *testing.T) {
		task, _ := ParseTask("Contents")

		AssertTaskText(t, task, "Contents")
		AssertTaskFinishedDate(t, task, time.Unix(0, 0))

		if task.creationDate.IsZero() {
			t.Fatal("task.creationDate has zero value")
		}
	})

	t.Run("escaped_no_separator", func(t *testing.T) {
		task, _ := ParseTask(`hello\|world`)

		AssertTaskText(t, task, "hello|world")
		AssertTaskFinishedDate(t, task, time.Unix(0, 0))
	})

	t.Run("has_default_creationDate", func(t *testing.T) {
		task, _ := ParseTask("hi | finished:2001/03/14/23/58")

		var dt time.Time
		AssertNotEqual(t, task.creationDate, dt, "Task's creation date is zero-value")
		AssertTaskFinishedDate(t, task, time.Date(2001, 03, 14, 23, 58, 0, 0, time.UTC))
	})

	t.Run("full", func(t *testing.T) {
		task, _ := ParseTask("Example Task | id: 7b91fb49a85ea06bb0276e70984d602e62e95ea5, creation:2003/04/15/22/18, finished:2001/01/01/00/00")

		AssertTaskText(t, task, "Example Task")
		AssertTaskHash(t, task, "7b91fb49a85ea06bb0276e70984d602e62e95ea5")
		AssertTaskCreationDate(t, task, time.Date(2003, 4, 15, 22, 18, 0, 0, time.UTC))
		AssertTaskFinishedDate(t, task, time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC))
	})

	t.Run("full_random_meta_order", func(t *testing.T) {
		task, _ := ParseTask(`Example\|Task | finished:2001/01/01/00/00,id: 7b91fb49a85ea06bb0276e70984d602e62e95ea5 ,creation:2003/04/15/22/18`)

		AssertTaskText(t, task, "Example|Task")
		AssertTaskHash(t, task, "7b91fb49a85ea06bb0276e70984d602e62e95ea5")
		AssertTaskCreationDate(t, task, time.Date(2003, time.April, 15, 22, 18, 0, 0, time.UTC))
		AssertTaskFinishedDate(t, task, time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC))
	})
}
