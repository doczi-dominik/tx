package main

import (
	"os"

	"github.com/spf13/afero"
)

type FS struct {
	Fs afero.Fs

	initialized  bool
	TaskfilePath string
	SyncfilePath string
	DonefilePath string
}

func createOsFS() *FS {
	return &FS{
		Fs: afero.NewOsFs(),
	}
}

func createMemFS() *FS {
	return &FS{
		Fs:           afero.NewMemMapFs(),
		TaskfilePath: "mock-tasks",
		DonefilePath: ".mock-tasks.done",
		SyncfilePath: ".mock-tasks.sync",
		initialized:  true,
	}
}

func (fs *FS) init() {
	if fs.initialized {
		return
	}

	tfPath := ConfigOptions.List

	fs.TaskfilePath = tfPath
	fs.DonefilePath = GetMetafilePath(".done", tfPath)
	fs.SyncfilePath = GetMetafilePath(".sync", tfPath)

	fs.initialized = true
}

func (fs *FS) getTaskfilePath() string {
	return fs.TaskfilePath
}

func (fs *FS) getDonefilePath() string {
	return fs.DonefilePath
}

func (fs *FS) getSyncfilePath() string {
	return fs.SyncfilePath
}

func (fs *FS) createFile(path string, errorCode int) afero.File {
	return createFile(fs.Fs, path, errorCode)
}

func (fs *FS) openTaskfile(optional bool) (taskfile afero.File) {
	return openFile(fs.Fs, fs.TaskfilePath, ErrTaskfileOpen, optional)
}

func (fs *FS) createTaskfile() (taskfile afero.File) {
	return createFile(fs.Fs, fs.TaskfilePath, ErrTaskfileOpen)
}

func (fs *FS) openSyncfile(optional bool) (syncfile afero.File) {
	return openFile(fs.Fs, fs.SyncfilePath, ErrSyncfileOpen, optional)
}

func (fs *FS) createSyncfile() (syncfile afero.File) {
	return createFile(fs.Fs, fs.SyncfilePath, ErrSyncfileOpen)
}

func (fs *FS) openSyncfileRW() (syncfile afero.File) {
	file, err := fs.Fs.OpenFile(fs.SyncfilePath, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		Error(ErrSyncfileOpen, fs.SyncfilePath, err)
		return nil
	}

	return file
}

func (fs *FS) Stat(name string) (os.FileInfo, error) {
	return fs.Fs.Stat(name)
}

func (fs *FS) Remove(name string) error {
	return fs.Fs.Remove(name)
}

func (fs *FS) Truncate(name string) error {
	file, err := fs.Fs.OpenFile(name, os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	return file.Close()
}
