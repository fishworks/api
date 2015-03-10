package api

import (
	"testing"
	"time"
)

func TestCreateApp(t *testing.T) {
	app := NewApp("test")

	// close enough for government work!
	if time.Since(app.Created) > time.Duration(1*time.Millisecond) {
		t.Errorf("Expected the app's created date to be within 1ms, got '%v'", time.Since(app.Created))
	}

	if time.Since(app.Updated) > time.Duration(1*time.Millisecond) {
		t.Errorf("Expected the app's updated date to be within 1ms, got '%v'", time.Since(app.Updated))
	}

	if app.UUID == "" {
		t.Error("UUID is uninitialized")
	}

	if app.ID != "test" {
		t.Errorf("expected app ID == '%s', got '%s'", "test", app.ID)
	}
}

func TestCreateAppWithNoID(t *testing.T) {
	app := NewApp("")

	if app.ID == "" {
		t.Error("expected app ID to be generated, got empty string")
	}
}
