package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var storage string

func operationHandler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")

		w.Write([]byte(`
	{
		"contents": "A | id:6dcd4ce23d88e2ee9568ba546c007c63d9131c1b, creation:2021/01/28/09/59, finished:1970/01/01/01/00\n",
		"doneContents": "B | id:ae4f281df5a5d0ff3cad6371f76d5c29b6d953ec, creation:2021/01/28/10/00, finished:1970/01/01/01/00\n"
	}
	`))
	case "PUT":
		var json []byte
		_, _ = req.Body.Read(json)
		req.Body.Close()

		storage = string(json)
		w.WriteHeader(http.StatusOK)
	case "DELETE":
		storage = ""

		w.WriteHeader(http.StatusOK)
	}
}

func newSyncIDHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		w.Header().Add("Location", "/testing-sync-id")

		w.WriteHeader(http.StatusCreated)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func TestSyncing(t *testing.T) {
	InitTestingFS(t)

	// Set up a dummy Sync service.
	mux := http.NewServeMux()

	mux.HandleFunc("/", newSyncIDHandler)
	mux.HandleFunc("/testing-sync-id", operationHandler)
	mux.HandleFunc("/another-sync-id", operationHandler)

	ts := httptest.NewServer(mux)
	switchingTs := httptest.NewServer(mux)

	{
		// Test `sync enable` with requested Sync ID
		params := EnableParams{
			URL: ts.URL,
			Args: struct {
				SyncID string "description:\"The Sync ID to pair with the tasklist. If unspecified, a new Sync ID will be requested from the provided Sync Service.\""
			}{},
		}

		params.Execute([]string{})

		lm := &TasklistManager{}
		lm.ParseSyncfile()

		AssertEqual(t, lm.syncID, "testing-sync-id", "Simulated TasklistManager's syncID field is not \"testing-sync-id\"")
	}
	{
		// Test `sync disable`
		params := DisableParams{}

		params.Execute([]string{})

		lm := &TasklistManager{}
		lm.ParseSyncfile()

		AssertEqual(t, lm.syncID, "", "Simulated TasklistManager's syncID field is not empty")
		AssertEqual(t, lm.syncURL, "", "Simulated TasklistManager's syncURL field is not empty")
	}
	{
		// Test `sync enable` with provided Sync ID
		params := EnableParams{
			URL: ts.URL,
			Args: struct {
				SyncID string "description:\"The Sync ID to pair with the tasklist. If unspecified, a new Sync ID will be requested from the provided Sync Service.\""
			}{},
		}

		params.Args.SyncID = "another-sync-id"
		params.Execute([]string{})

		lm := &TasklistManager{}
		lm.ParseSyncfile()

		AssertEqual(t, lm.syncID, "another-sync-id", "Simulated TasklistManager's syncID field is not \"another-sync-id\"")
	}
	{
		// Test `sync free`
		storage = "DUMMY"

		params := FreeParams{
			Args: struct {
				SyncID string "description:\"The Sync ID of the desired tasklist. If unspecified, the Sync ID will be read from the appropriate Syncfile.\""
			}{},
		}

		params.Execute([]string{})

		AssertEqual(t, storage, "", "Dummy storage is not empty after DELETE Request")

		lm := &TasklistManager{}
		lm.ParseSyncfile()

		AssertEqual(t, lm.syncID, "", "Simulated TasklistManager's syncID field is not empty")
		AssertEqual(t, lm.syncURL, "", "Simulated TasklistManager's syncURL field is not empty")
	}
	{
		// Test `sync switch`
		enableParams := EnableParams{
			URL: ts.URL,
			Args: struct {
				SyncID string "description:\"The Sync ID to pair with the tasklist. If unspecified, a new Sync ID will be requested from the provided Sync Service.\""
			}{},
		}

		enableParams.Execute([]string{})

		params := SwitchParams{
			Args: struct {
				URL    string "description:\"The URL of the Sync service. If unspecified, the fallback Sync URL will be used.\""
				SyncID string "description:\"The new Sync ID of the tasklist. If unspecified, a new Sync ID will be requested.\""
			}{},
		}

		params.Args.URL = switchingTs.URL
		params.Execute([]string{})

		lm := &TasklistManager{}
		lm.ParseSyncfile()

		AssertEqual(t, lm.syncURL, switchingTs.URL+"/", "Simulated TasklistManager's syncURL field has not changed")
	}
}
