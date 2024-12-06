package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// JSONHeaders are the set of HTTP headers required to accept the body in JSON.
var JSONHeaders = http.Header{
	"Content-Type": {"application/json"},
	"Accept":       {"application/json"},
}

// LoadSource is a type used for the Local and Network constants.
type LoadSource int

const (
	// Unavailable notes that the tasklist manager is still loading.
	Unavailable LoadSource = iota
	// Local notes that the tasklist manager opted to use local files.
	Local
	// Network notes that the tasklist manager opted to load from a sync service.
	Network
	// OutdatedNetwork notes that the tasklist manager detected that the local
	// taskfile is newer than the last network update, and will upload the
	// local tasklist to the appropriate sync service.
	OutdatedNetwork
)

// TasklistManager is responsible for loading and saving taskfiles, both local
// files and Sync service data.
type TasklistManager struct {
	source LoadSource

	syncID            string
	syncURL           string
	lastNetworkUpdate time.Time
}

// ParseSyncfile reads the appropriate syncfile and sets the manager's fields
// to their values.
func (tm *TasklistManager) ParseSyncfile() {
	syncfile := GlobalFS.openSyncfile(true)

	if syncfile == nil {
		return
	}

	defer syncfile.Close()

	scanner := bufio.NewScanner(syncfile)

	for scanner.Scan() {
		line := scanner.Text()

		// Extract Sync ID
		syncIDMatch := SyncIDPattern.FindStringSubmatch(line)

		if len(syncIDMatch) == 2 {
			tm.syncID = syncIDMatch[1]
			continue
		}

		// Extract Sync URL
		syncURLMatch := SyncURLPattern.FindStringSubmatch(line)

		if len(syncURLMatch) == 2 {
			tm.syncURL = syncURLMatch[1]
			continue
		}

		// Extract lastNetworkUpdate
		lastNetworkUpdateMatch := LastNetworkUpdatePattern.FindStringSubmatch(line)

		if len(lastNetworkUpdateMatch) == 2 {
			lastNetworkUpdate, err := time.Parse(LastNetworkUpdateFormat, lastNetworkUpdateMatch[1])

			if err != nil {
				Warn("Invalid lastNetworkUpdate date in syncfile \"%s\": %s", GlobalFS.getSyncfilePath(), lastNetworkUpdateMatch[1])
			}

			tm.lastNetworkUpdate = lastNetworkUpdate
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		Error(ErrSyncfileRead, GlobalFS.getSyncfilePath(), err)
	}
}

// Init sets all the filepaths and initializes both active and finished
// tasklists.
func (tm *TasklistManager) Init() {
	GlobalFS.init()

	MainList.filePath = GlobalFS.getTaskfilePath()
	DoneList.filePath = GlobalFS.getDonefilePath()

	MainList.tasks = make(map[int]Task)
	DoneList.tasks = make(map[int]Task)
}

func (tm *TasklistManager) loadNetwork() LoadSource {
	// Prepare URL
	tm.syncURL = strings.TrimSpace(tm.syncURL)
	EnsureTrailingSlash(&tm.syncURL)

	tm.syncID = strings.TrimSpace(tm.syncID)

	url := tm.syncURL + tm.syncID

	// Complete GET Request
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		Error(ErrGETReqCreate, url, err)
	}

	request.Header = JSONHeaders

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		Warn("Could not complete GET request: %v", err)
		return Local
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode

	switch statusCode {
	case 404:
		Warn("Could not find tasklist with Sync ID \"%s\"", tm.syncID)
		return Local
	case 200:
		// Unmarshal JSON stirng into map[string]string
		var jsonData []byte

		for s := bufio.NewScanner(resp.Body); s.Scan(); {
			jsonData = append(jsonData, s.Bytes()...)
		}

		var data map[string]interface{}

		if json.Unmarshal(jsonData, &data) != nil {
			Warn("Failed to parse network response: %v", err)
			return Local
		}

		// Parse active and finished tasklists
		{
			reader := strings.NewReader(data["contents"].(string))

			MainList.ParseTasklines("[syncID:"+tm.syncID+"]", reader)
		}

		{
			reader := strings.NewReader(data["doneContents"].(string))

			DoneList.ParseTasklines("[syncID:"+tm.syncID+"]", reader)
		}

		MainList.loaded = true
		DoneList.loaded = true

		return Network

	default:
		Warn("Invalid response when loading tasklist: status is " + resp.Status)
		return Local
	}
}

// Load loads from network if possible and determines the loading source.
func (tm *TasklistManager) Load() {
	tm.source = func() LoadSource {
		// "Bouncer" statements
		if ConfigOptions.Offline {
			return Local
		}

		tm.ParseSyncfile()

		if tm.syncID == "" {
			return Local
		}

		if tm.syncURL == "" {
			tm.syncURL = ConfigOptions.FallbackSyncURL
		}

		if MainList.MTimeAfter(tm.lastNetworkUpdate) || DoneList.MTimeAfter(tm.lastNetworkUpdate) {
			Warn("Local tasklist is newer than the synced one. Uploading local tasklist.")
			return OutdatedNetwork
		}

		return tm.loadNetwork()
	}()
}

// EnsureInitialized makes sure that the tasklist is loaded and does so if not.
func (tm *TasklistManager) EnsureInitialized(tasklist *Tasklist) {
	if tasklist.loaded {
		return
	}

	if tm.source == Unavailable {
		tm.Init()
		tm.Load()
	}

	if tm.source == Local || tm.source == OutdatedNetwork {
		tasklist.LoadLocal()

		if tm.source == OutdatedNetwork {
			tasklist.MarkModified()
		}
	}
}

// Save is responsible for saving changes to the local taskfile and uploading
// data to the appropriate Sync service.
func (tm *TasklistManager) Save() {
	if !MainList.modified && !DoneList.modified {
		return
	}

	MainList.SerializeTasks()
	DoneList.SerializeTasks()

	// Save local taskfiles
	if MainList.modified {
		MainList.SaveLocal()
	}

	if DoneList.modified {
		DoneList.SaveLocal()
	}

	// Upload to Sync service
	if tm.source > Local {
		// Marshal into JSON
		jsonMap := map[string]string{
			"contents":     string(MainList.serialized),
			"doneContents": string(DoneList.serialized),
		}

		jsonData, err := json.Marshal(jsonMap)

		if err != nil {
			Error(ErrSerializeJSON, err)
		}

		// Complete PUT Request
		url := tm.syncURL + tm.syncID
		request, err := http.NewRequest("PUT", url, bytes.NewReader(jsonData))

		if err != nil {
			Error(ErrPUTReqCreate, url, err)
		}

		request.Header = JSONHeaders

		resp, err := http.DefaultClient.Do(request)

		if err != nil {
			Error(ErrPUTReqComplete, url, err)
		}

		defer resp.Body.Close()

		statusCode := resp.StatusCode

		switch statusCode {
		case 404:
			Error(ErrNoSyncID, tm.syncID)
		case 200:
			// Update LastNetworkUpdate value in Syncfile
			newLastNetworkUpdate := StripNanoFromTime(time.Now())
			line := "lastNetworkUpdate: " + newLastNetworkUpdate.Format(LastNetworkUpdateFormat) + "\n"

			ReplaceOrAppendSyncfileLine(LastNetworkUpdatePattern, line)
		default:
			Error(ErrInvalidResponse, url, "status is "+resp.Status)
		}
	}

	RunCallback()
}
