package main

import "testing"

func TestRunCallback(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		AssertExitSuccess(t, "TestRunCallback/success", func() {
			ConfigOptions.Callback = "echo ab | sed s/a/hello/ | sed 's/b/ world'"
			RunCallback()
		})
	})

	t.Run("failure", func(t *testing.T) {
		AssertExitError(t, "TestRunCallback/failure", ErrCallback, func() {
			ConfigOptions.Callback = "false"
			RunCallback()
		})
	})
}
