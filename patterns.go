package main

import "regexp"

// SyncIDPattern is used for extracting the Sync ID from a syncfile.
var SyncIDPattern = regexp.MustCompile(`(?i)syncid[ \t]*:[ \t]*(.+)`)

// SyncURLPattern is used for extracting the Sync service URL from a syncfile.
var SyncURLPattern = regexp.MustCompile(`(?i)syncurl[ \t]*:[ \t]*(.+)`)

// LastNetworkUpdatePattern is used for extracting the last successful HTTP request
// date from a syncfile.
var LastNetworkUpdatePattern = regexp.MustCompile(`(?i)lastnetworkupdate[ \t]*:[ \t](\d{4}/\d{2}/\d{2}/\d{2}\/\d{2}\/\d{2})`)

// SeparatorPattern is used for finding the text/meta separator pipe in a
// taskline.
var SeparatorPattern = regexp.MustCompile(`[^\\]\|`)

// CreationDatePattern is used for extracting the creation date from a
// taskline.
var CreationDatePattern = regexp.MustCompile(`(?i)creation:[\t ]*(\d{4}/\d{2}/\d{2}/\d{2}\/\d{2})`)

// FinishedDatePattern is used for extracting the date the task was marked as
// finished from a taskline.
var FinishedDatePattern = regexp.MustCompile(`(?i)finished:[\t ]*(\d{4}/\d{2}/\d{2}/\d{2}\/\d{2})`)

// HashPattern is used for extracting the SHA-1 sum of the task's text from a
// taskline.
var HashPattern = regexp.MustCompile(`(?i)id:[\t ]*([a-f0-9]{40})`)

// FullDateFormat specifies the general format for parsing strings to time
// objects.
var FullDateFormat = "2006/01/02/15/04"

// DateFormat specifies the format for poarsing the date portion (Y, M, D) of
// stirngs to time objects.
var DateFormat = FullDateFormat[:10]

// LastNetworkUpdateFormat specifies the format for parsing the
// lastNetworkUpdate value into a time object.
var LastNetworkUpdateFormat = "2006/01/02/15/04/05"

// DisplayTimeFormat specifies how a time object's time portion (H, M) should
// look when displayed in a human-readable form.
var DisplayTimeFormat = "15:04"
