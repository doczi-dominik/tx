package main

import (
	"bufio"
	"os"
	"regexp"
)

func openFile(filePath string, errorCode int, optional bool) *os.File {
	file, err := os.Open(filePath)

	if err != nil {
		if optional && os.IsNotExist(err) {
			return nil
		}

		Error(errorCode, filePath, err)
	}

	return file
}

func createFile(filePath string, errorCode int) *os.File {
	file, err := os.Create(filePath)

	if err != nil {
		Error(errorCode, filePath, err)
	}

	return file
}

// ReplaceOrAppendSyncfileLine essentialy makes sure that an "entry" appears in a
// syncfile exactly once.
func ReplaceOrAppendSyncfileLine(pattern *regexp.Regexp, repl string) {
	// Read Syncfile contents into an array, skipping lines that match the
	// pattern.
	syncfile := GlobalFS.openSyncfileRW()

	defer syncfile.Close()

	var contents []byte

	scanner := bufio.NewScanner(syncfile)

	for scanner.Scan() {
		line := scanner.Text()

		if !pattern.MatchString(line) {
			contents = append(contents, line+"\n"...)
		}
	}

	if err := scanner.Err(); err != nil {
		Error(ErrSyncfileRead, GlobalFS.getSyncfilePath(), err)
	}

	// Add the replacement to the end of the list and write the array back into
	// the file.
	contents = append(contents, repl...)

	err := syncfile.Truncate(0)

	if err != nil {
		Error(ErrSyncfileWrite, GlobalFS.getSyncfilePath(), err)
	}

	syncfile.Seek(0, 0)
	_, err = syncfile.Write(contents)

	if err != nil {
		Error(ErrSyncfileWrite, GlobalFS.getSyncfilePath(), err)
	}
}
