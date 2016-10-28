package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fishworks/api"
)

func clearDB() {
	Apps = Apps[:0]
}

func TestEmptyListAppsReturnsNoContent(t *testing.T) {
	srv, err := New("tcp", "0.0.0.0:4567")
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()
	r := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/apps", nil)
	if err != nil {
		t.Fatal(err)
	}
	srv.ServeRequest(r, req)
	if r.Code != http.StatusNoContent {
		t.Fatalf("%d NO CONTENT expected, received %d\n", http.StatusNoContent, r.Code)
	}
}

func TestCreateAppAndThenList(t *testing.T) {
	defer clearDB()
	srv, err := New("tcp", "0.0.0.0:4567")
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
	if len(Apps) != 1 {
		t.Fatalf("%d app expected, got %d", 1, len(Apps))
	}
}

func TestCreateAppWithID(t *testing.T) {
	defer clearDB()
	srv, err := New("tcp", "0.0.0.0:4567")
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()
	r := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/apps", bytes.NewBuffer([]byte(`{"id":"autotest"}`)))
	if err != nil {
		t.Fatal(err)
	}
	srv.ServeRequest(r, req)
	if r.Code != http.StatusCreated {
		t.Fatalf("%d CREATED expected, received %d\n", http.StatusCreated, r.Code)
	}
	if Apps[0].ID != "autotest" {
		t.Errorf("%s expected, received %s\n", "autotest", Apps[0].ID)
	}
}

// TestGetAppRemovesUUID tests that an application's UUID does not show up in the response body.
func TestGetAppRemovesUUID(t *testing.T) {
	defer clearDB()
	app := api.NewApp("autotest")
	Apps = append(Apps, app)
	srv, err := New("tcp", "0.0.0.0:4567")
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()
	r := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/apps/autotest", nil)
	if err != nil {
		t.Fatal(err)
	}
	srv.ServeRequest(r, req)
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

func TestGetAppLogs(t *testing.T) {
	defer clearDB()
	app := api.NewApp("autotest")
	Apps = append(Apps, app)
	app.Log("ohai der =3")
	srv, err := New("tcp", "0.0.0.0:4567")
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()
	r := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/apps/autotest/logs", nil)
	if err != nil {
		t.Fatal(err)
	}
	srv.ServeRequest(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("%d OK expected, received %d\n", http.StatusOK, r.Code)
	}
	if !strings.Contains(r.Body.String(), "deis[api]: ohai der =3") {
		t.Fatalf("%s expected, received %s\n", "deis[api]: ohai der =3", r.Body.String())
	}
}

func TestDeleteApp(t *testing.T) {
	defer clearDB()
	app := api.NewApp("autotest")
	Apps = append(Apps, app)
	srv, err := New("tcp", "0.0.0.0:4567")
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()
	r := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/apps/autotest", nil)
	if err != nil {
		t.Fatal(err)
	}
	srv.ServeRequest(r, req)
	if r.Code != http.StatusNoContent {
		t.Fatalf("%d NO CONTENT expected, received %d\n", http.StatusNoContent, r.Code)
	}
	if len(Apps) != 0 {
		t.Fatalf("%d expected, received %d\n", 0, len(Apps))
	}
}
