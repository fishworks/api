package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fishworks/api"
	"github.com/fishworks/api/server"
)

func TestEmptyListAppsReturnsNoContent(t *testing.T) {
	server, err := server.NewServer("tcp", "0.0.0.0:4567")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	r := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/apps", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeRequest(r, req)
	if r.Code != http.StatusNoContent {
		t.Fatalf("%d NO CONTENT expected, received %d\n", http.StatusNoContent, r.Code)
	}
}

func TestCreateAppAndThenList(t *testing.T) {
	srv, err := server.NewServer("tcp", "0.0.0.0:4567")
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()
	r := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/apps", nil)
	if err != nil {
		t.Fatal(err)
	}
	srv.ServeRequest(r, req)
	if r.Code != http.StatusCreated {
		t.Fatalf("%d CREATED expected, received %d\n", http.StatusCreated, r.Code)
	}
	r = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/apps", nil)
	if err != nil {
		t.Fatal(err)
	}
	srv.ServeRequest(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("%d OK expected, received %d\n", http.StatusOK, r.Code)
	}
	if len(server.Apps) != 1 {
		t.Fatalf("%d app expected, got %d", 1, len(server.Apps))
	}
}

// TestGetAppRemovesUUID tests that an application's UUID does not show up in the response body.
func TestGetAppRemovesUUID(t *testing.T) {
	app := api.NewApp("autotest")
	server.Apps = append(server.Apps, app)
	server, err := server.NewServer("tcp", "0.0.0.0:4567")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	r := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/apps/autotest", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeRequest(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("%d OK expected, received %d\n", http.StatusOK, r.Code)
	}
	var app2 api.App
	if err := json.Unmarshal(r.Body.Bytes(), &app2); err != nil {
		t.Fatal(err)
	}
	if app2.UUID != "" {
		t.Fatalf("%s expected, received %s\n", "blank UUID", app2.UUID)
	}
}
