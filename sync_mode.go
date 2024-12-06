package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// SyncActions contains all options and positional arguments for syncfile
// management mode.
type SyncActions struct {
	Enable  EnableParams  `command:"enable" description:"Enable the syncing of tasklists using a Sync service."`
	Disable DisableParams `command:"disable" description:"Remove the Sync ID (if any) from the syncfile"`
	Free    FreeParams    `command:"free" description:"Permanently deactivate a Sync ID. The tasklist it belongs to will be deleted from the Sync service."`
	Now     NowParams     `command:"now" description:"Uploads the local tasklist to the configured Sync service, even if the file is outdated."`
	Switch  SwitchParams  `command:"switch" description:"Change the Sync service to be used for this tasklist."`
}

// EnableParams holds the command line arguments for the `sync enable` subcommand.
type EnableParams struct {
	URL  string `short:"u" long:"sync-url" description:"The Sync Service to request the Sync ID from. If unspecified, the global fallback Sync URL will be used."`
	Args struct {
		SyncID string `description:"The Sync ID to pair with the tasklist. If unspecified, a new Sync ID will be requested from the provided Sync Service."`
	} `positional-args:"true"`
}

// Execute uses the provided EnableParams and enables syncing for the tasklist.
func (a *EnableParams) Execute(args []string) error {
	GlobalFS.init()

	var syncURL string

	if a.URL != "" {
		syncURL = strings.TrimSpace(a.URL)
	} else {
		syncURL = strings.TrimSpace(ConfigOptions.FallbackSyncURL)
	}

	EnsureTrailingSlash(&syncURL)

	enable(syncURL, a.Args.SyncID)

	return nil
}

// DisableParams holds the command line arguments for the `sync disable` subcommand.
type DisableParams struct{}

// Execute uses the provided DisableParams and disables syncing for the
// tasklist.
func (a *DisableParams) Execute(args []string) error {
	GlobalFS.init()

	ReplaceOrAppendSyncfileLine(SyncIDPattern, "")
	ReplaceOrAppendSyncfileLine(SyncURLPattern, "")

	return nil
}

// FreeParams holds the command line arguments for the `sync free` subcommand.
type FreeParams struct {
	Args struct {
		SyncID string `description:"The Sync ID of the desired tasklist. If unspecified, the Sync ID will be read from the appropriate Syncfile."`
	} `positional-args:"true"`
}

// Execute uses the provided FreeParams and "frees" storage on the sync service
// by deleting the uploaded tasklist and disabling syncing for the local one.
func (a *FreeParams) Execute(args []string) error {
	GlobalFS.init()

	// Read Sync ID and Sync URL from syncfile
	var (
		syncID           string
		syncURL          string
		syncfileContents string
	)

	if a.Args.SyncID != "" {
		syncID = a.Args.SyncID
	}

	syncfile := GlobalFS.openSyncfile(false)
	defer syncfile.Close()

	scanner := bufio.NewScanner(syncfile)

	for scanner.Scan() {
		line := scanner.Text()

		if syncID == "" {
			if m := SyncIDPattern.FindStringSubmatch(line); len(m) == 2 {
				syncID = strings.TrimSpace(m[1])
				continue
			}
		}

		if m := SyncURLPattern.FindStringSubmatch(line); len(m) == 2 {
			syncURL = strings.TrimSpace(m[1])
			continue
		}

		syncfileContents += line
	}

	if err := scanner.Err(); err != nil {
		Error(ErrSyncfileRead, GlobalFS.getSyncfilePath(), err)
	}

	if syncURL == "" {
		syncURL = strings.TrimSpace(ConfigOptions.FallbackSyncURL)
	}

	EnsureTrailingSlash(&syncURL)

	// Request deletion from the Sync Service
	url := syncURL + syncID
	request, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		Error(ErrDELETEReqCreate, url, err)
	}

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		Error(ErrDELETEReqComplete, url, err)
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode

	switch statusCode {
	case 404:
		Error(ErrNoSyncID, syncID)
	case 405:
		Error(ErrUnsupportedConfig, url)
	case 200:
		// Remove the syncID line from the syncfile
		syncfile = GlobalFS.createSyncfile()
		defer syncfile.Close()

		_, err = syncfile.WriteString(syncfileContents + "\n")

		if err != nil {
			Error(ErrSyncfileWrite, GlobalFS.getSyncfilePath(), err)
		}
	default:
		Error(ErrInvalidResponse, url, "status is "+resp.Status)
	}

	return nil
}

// NowParams holds the command line arguments for the `sync now` subcommand.
type NowParams struct{}

// Execute uses the provided NowParams and uploads the local tasklist to the
// Sync service.
func (a *NowParams) Execute(args []string) error {
	upload()

	return nil
}

// SwitchParams holds the command line arguments for the `sync switch` subcommand.
type SwitchParams struct {
	Args struct {
		URL    string `description:"The URL of the Sync service. If unspecified, the fallback Sync URL will be used."`
		SyncID string `description:"The new Sync ID of the tasklist. If unspecified, a new Sync ID will be requested."`
	} `positional-args:"true"`
}

// Execute uses the provided SwitchParams and switches Sync services for the
// current tasklist.
func (a *SwitchParams) Execute(args []string) error {
	GlobalFS.init()

	var syncURL string

	if a.Args.URL == "" {
		syncURL = strings.TrimSpace(ConfigOptions.FallbackSyncURL)
	} else {
		syncURL = strings.TrimSpace(a.Args.URL)
	}

	EnsureTrailingSlash(&syncURL)

	enable(syncURL, a.Args.SyncID)

	return nil
}

// init gets called when the package is imported; assigns functions to the
// respective action and adds the subcommand to the global argument parser.
func init() {
	var actions SyncActions
	GlobalParser.AddCommand("sync", "Manage Syncing for the current tasklist", "", &actions)
}

func enable(syncURL string, syncID string) {
	var newSyncID string

	if syncID != "" {
		newSyncID = syncID
	} else {
		resp, err := http.Post(syncURL, "application/json", nil)

		if err != nil {
			Error(ErrRequestSyncID, syncURL, err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != 201 {
			Error(ErrInvalidResponse, syncURL, "status is "+resp.Status)
		}

		// Extract the new syncID from the updated location header.
		location, err := resp.Location()

		if err != nil {
			Error(ErrInvalidResponse, err)
		}

		locationString := location.String()
		newSyncID = locationString[strings.LastIndex(locationString, "/")+1:]
	}

	newSyncID = strings.TrimSpace(newSyncID)

	ReplaceOrAppendSyncfileLine(SyncIDPattern, "syncID: "+newSyncID+"\n")
	ReplaceOrAppendSyncfileLine(SyncURLPattern, "syncURL: "+syncURL+"\n")

	upload()

	fmt.Printf("\"%s\" from \"%s\"", newSyncID, syncURL)
}

func upload() {
	// Simulate an outdated taskfile scenario using a TasklistMangager and
	// upload.
	lm := &TasklistManager{
		source: OutdatedNetwork,
	}

	lm.Init()

	lm.ParseSyncfile()

	MainList.LoadLocal()
	DoneList.LoadLocal()

	MainList.MarkModified()
	DoneList.MarkModified()

	lm.Save()
}
