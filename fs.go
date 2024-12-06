package main

import "os"

type FS interface {
	init()

	getSyncfilePath() string
	getTaskfilePath() string
	getDonefilePath() string

	createFile(path string, errorCode int) (file *os.File)

	openTaskfile(optional bool) (taskfile *os.File)
	createTaskfile() (taskfile *os.File)
	openSyncfile(optional bool) (syncfile *os.File)
	createSyncfile() (syncfile *os.File)
	openSyncfileRW() (syncfile *os.File)
}
