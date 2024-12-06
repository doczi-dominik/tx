package main

import "os"

type FileFS struct {
	initialized  bool
	TaskfilePath string
	SyncfilePath string
	DonefilePath string
}

func createFileFS() *FileFS {
	return &FileFS{
		initialized: false,
	}
}

func (fs *FileFS) init() {
	if fs.initialized {
		return
	}

	tfPath := ConfigOptions.List

	fs.TaskfilePath = tfPath
	fs.DonefilePath = GetMetafilePath(".done", tfPath)
	fs.SyncfilePath = GetMetafilePath(".sync", tfPath)

	fs.initialized = true
}

func (fs *FileFS) getTaskfilePath() string {
	return fs.TaskfilePath
}

func (fs *FileFS) getDonefilePath() string {
	return fs.DonefilePath
}

func (fs *FileFS) getSyncfilePath() string {
	return fs.SyncfilePath
}

func (fs *FileFS) createFile(path string, errorCode int) *os.File {
	return createFile(path, errorCode)
}

func (fs *FileFS) openTaskfile(optional bool) (taskfile *os.File) {
	return openFile(fs.TaskfilePath, ErrTaskfileOpen, optional)
}

func (fs *FileFS) createTaskfile() (taskfile *os.File) {
	return createFile(fs.TaskfilePath, ErrTaskfileOpen)
}

func (fs *FileFS) openSyncfile(optional bool) (syncfile *os.File) {
	return openFile(fs.SyncfilePath, ErrSyncfileOpen, optional)
}

func (fs *FileFS) createSyncfile() (syncfile *os.File) {
	return createFile(fs.SyncfilePath, ErrSyncfileOpen)
}

func (fs *FileFS) openSyncfileRW() (syncfile *os.File) {
	file, err := os.OpenFile(fs.SyncfilePath, os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		Error(ErrSyncfileOpen, fs.SyncfilePath, err)
		return nil
	}

	return file
}
