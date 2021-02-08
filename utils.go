package main

import (
	"path"
	"strings"
	"time"
)

// GetMetafilePath derives the path of a metafile from the regular tasklist's
// path.
func GetMetafilePath(ext string) string {
	dir, filename := path.Split(TaskfilePath)

	return dir + "." + filename + ext
}

// StripNanoFromTime zeroes the nanosecond field of a time object.
func StripNanoFromTime(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		0,
		time.UTC,
	)
}

// EnsureTrailingSlash appends a slash ("/"") to a string if it does end with
// one (plus a newline).
func EnsureTrailingSlash(s *string) {
	if !strings.HasSuffix(*s, "/") {
		*s += "/"
	}
}
