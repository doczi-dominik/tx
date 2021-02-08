package main

import (
	"fmt"
	"os"
)

const (
	// ErrFlagParsing is used when the global argument parser returns an error.
	// Message requires an error (type error).
	ErrFlagParsing = iota
	// ErrInvalidSelector is used when an invalid selector is passed to a
	// tasklist operation. Message requires name of the enclosing operation
	// (type string) the invalid selector (type string).
	ErrInvalidSelector
	// ErrEditInvalidSelector is a variant of ErrInvalidSelector for the Edit
	// action. Message requires the invalid selector (type string).
	ErrEditInvalidSelector
	// ErrInvalidIndex is used when a provided selector refers to a
	// non-existant task. Message requires the name of the enclosing operation
	// (type string) and the invalid index (type int).
	ErrInvalidIndex
	// ErrTaskValidation is used when a task is deemed to be invalid
	// (e.g.: text contains a newline). Message requires the name of the
	// enclosing operation (type string) and the error (type error).
	ErrTaskValidation
	// ErrTasklistEmpty is used when a non-creative operation is used on an
	// empty tasklist. Message requires the name of the enclosing operation
	// (type string).
	ErrTasklistEmpty
	// ErrSerializeJSON is used when a tasklist cannot be converted to JSON.
	// Message requires an error (type error).
	ErrSerializeJSON
	// ErrCallback is used when the specifed callback shell command
	// cannot execute or it returns a non-zero error code. Message requires
	// the callback command (type string) and the error (type error).
	ErrCallback
)

// ErrTaskfile[...] are used when handling file operations on a regular
// taskfile. Messages require the filepath (type string) and the error (type error).
const (
	ErrTaskfileOpen = 8 + iota
	ErrTaskfileWrite
	ErrTaskfileRead
)

// ErrBackup[...] are used when creating/writing a backup of a taskfile.
// Messages require the filepath (type string) and the error (type error).
const (
	ErrBackupCreate = 11 + iota
	ErrBackupWrite
)

// ErrSyncfile[...] are used when handling file operations on a syncfile.
// Messages require the filepath (type string) and the error (type error).
const (
	ErrSyncfileOpen = 13 + iota
	ErrSyncfileWrite
	ErrSyncfileRead
)

// ErrGETReq[...] are used when creating or completing a HTTP GET request.
// Messages require the URL (type string) and the error (type error).
const (
	ErrGETReqCreate = 16 + iota
)

// ErrPUTReq[...] are used when creating or completing a HTTP PUT request.
// Messages require the URL (type string) and the error (type error).
const (
	ErrPUTReqCreate = 17 + iota
	ErrPUTReqComplete
)

// ErrDELETEReq[...] are used when creating or completing a HTTP DELETE
// request. Messages require the URL (type string) and the error (type error).
const (
	ErrDELETEReqCreate = 19 + iota
	ErrDELETEReqComplete
)

const (
	// ErrNoSyncID indicates that no task exists with the given Sync ID.
	// Message requires the Sync ID (type string).
	ErrNoSyncID = 21 + iota
	// ErrRequestSyncID indicates that tx failed to request a new Sync ID from
	// the sync service. Message requires the URL (type string) and the error
	// (type error).
	ErrRequestSyncID
	// ErrInvalidResponse is used when an unparseable response is received from
	// the sync service. Message requires the URL (type string) and the error
	// (type error).
	ErrInvalidResponse
	// ErrUnsupportedConfig is used when trying to delete a tasklist from
	// the sync service while deletion is disabled on the sync service - an
	// operation which tx will never complete. Message requires the URL
	// (type string).
	ErrUnsupportedConfig
)

var errorMessages = [25]string{
	"Argument parser: %v",
	"%s: Invalid selector: \"%s\". Use --help for selector format information.",
	"Edit: Invalid selector: \"%s\". Use SELECT/NEW or SELECT/OLD/NEW.",
	"%s: Invalid index: %d",
	"%s: Task validation failed: %v",
	"%s: Tasklist is empty",
	"Could not serialize tasklist to JSON: %v",
	"Could not run callback \"%s\": %v",

	"Could not open taskfile \"%s\": %v",
	"Could not write taskfile \"%s\": %v",
	"Could not read taskfile \"%s\": %v",

	"Could not create backup file \"%s\": %v",
	"Could not write backup file \"%s\": %v",

	"Could not open syncfile \"%s\": %v",
	"Could not write syncfile \"%s\": %v",
	"Could not read syncfile \"%s\": %v",

	"Could not create GET request to \"%s\": %v",

	"Could not create PUT request to \"%s\": %v",
	"Could not complete PUT request for \"%s\": %v",

	"Could not create DELETE request to \"%s\": %v",
	"Could not complete DELETE request for \"%s\": %v",

	"No tasklist with Sync ID \"%s\"",
	"Could not request new Sync ID from \"%s\": %v",
	"Invalid response from \"%s\": %v",
	"Unsupported configuration: Deleting this tasklist is disabled on \"%s\"",
}

// Error is used to print a standard error message then exit.
func Error(code int, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "E: "+errorMessages[code]+"\n", args...)
	os.Exit(code + 1)
}

// Warn is used to print a standard warning message that does not cause tx to
// exit.
func Warn(message string, args ...interface{}) {
	if ConfigOptions.Quiet {
		return
	}

	fmt.Fprintf(os.Stderr, "W: "+message+"\n", args...)
}
