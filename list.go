package main

import (
	"io"
	"os"
	"time"
)

// Tasklist holds the loaded tasks and has methods to modify those tasks.
type Tasklist struct {
	filePath   string       // The path to the local taskfile.
	tasks      map[int]Task // Tasks mapped to an index.
	modified   bool         // States whether the tasklist has been modified.
	serialized []byte       // Stores the serialzed tasks before saving.
	loaded     bool         // True if the task has finished loading.
}

func (tl *Tasklist) copyToBackup() {
	src := GlobalFS.openTaskfile(true)

	if src != nil {
		defer src.Close()

		backupFilePath := GetMetafilePath(".bak", tl.filePath)
		dst := GlobalFS.createFile(backupFilePath, ErrBackupCreate)

		defer dst.Close()

		_, err := io.Copy(dst, src)

		if err != nil {
			Error(ErrBackupWrite, backupFilePath, err)
		}
	}
}

// LoadLocal reads the provided taskfile and parses tasks into the tasklist.
func (tl *Tasklist) LoadLocal() {
	taskfile := GlobalFS.openTaskfile(true)

	if taskfile != nil {
		defer taskfile.Close()

		tl.ParseTasklines(tl.filePath, taskfile)
	}

	tl.loaded = true
}

// SaveLocal serializes and writes tasks to the provided tasklist.
func (tl *Tasklist) SaveLocal() {
	if !ConfigOptions.Reckless {
		tl.copyToBackup()
	}

	if tl.IsEmpty() {
		if ConfigOptions.DeleteIfEmpty {
			err := GlobalFS.Remove(tl.filePath)

			if err != nil {
				Warn("Could not delete empty taskfile \"%s\": %v\n", tl.filePath, err)
			}
		} else {
			err := GlobalFS.Truncate(tl.filePath)

			if err != nil {
				Error(ErrTaskfileWrite, tl.filePath, err)
			}
		}
		return
	}

	taskfile := GlobalFS.createTaskfile()
	defer taskfile.Close()

	_, err := taskfile.Write(tl.serialized)

	if err != nil {
		Error(ErrTaskfileWrite, tl.filePath, err)
	}
}

// MTimeAfter determines if the last modification time for a taskfile is after
// another time object.
func (tl *Tasklist) MTimeAfter(compare time.Time) bool {
	stat, err := GlobalFS.Stat(tl.filePath)

	if err != nil {
		if !os.IsNotExist(err) {
			Warn("Could not retrieve mtime for taskfile \"%s\": %v", tl.filePath, err)
		}
		return false
	}

	return StripNanoFromTime(stat.ModTime()).After(compare)
}
