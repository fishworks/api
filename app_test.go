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
	if app.String() != "test" {
		t.Errorf("expected app.String() == '%s', got '%s'", "test", app.String())
	}
}

func TestCreateAppWithNoID(t *testing.T) {
	app := NewApp("")

	if app.ID == "" {
		t.Error("expected app ID to be generated, got empty string")
	}
}

func TestAppRelease(t *testing.T) {
	app := NewApp("")
	release := app.NewRelease(&Build{}, &Config{})
	if release == nil {
		t.Errorf("expected app to create a new release")
	}
	if release.Version != 1 {
		t.Errorf("expected version to be 1; got %d", release.Version)
	}
	if app.Ledger.Len() != 1 {
		t.Errorf("expected release to be appended to the ledger; got %d", app.Ledger.Len())
	}
	release2 := app.NewRelease(&Build{}, &Config{})
	if app.Ledger.Len() != 2 {
		t.Errorf("expected release2 to be appended to the ledger; got %d", app.Ledger.Len())
	}
	if !app.Ledger.Less(0, 1) {
		t.Errorf("expected v1 to be less than v2")
	}
	app.Ledger.Swap(0, 1)
	if app.Ledger[0] != release2 {
		t.Errorf("expected v2 to be first; got %s", app.Ledger[0])
	}
}

func TestAppRollback(t *testing.T) {
	app := NewApp("")
	release1 := app.NewRelease(&Build{}, &Config{})
	app.NewRelease(&Build{}, &Config{})

	// first, check that we cannot roll back to an invalid version
	if err := app.Rollback(0); err == nil {
		t.Errorf("expected rolling back to an invalid version number to error")
	}

	// now check that we cannot roll forward to a release that does not exist yet
	if err := app.Rollback(3); err == nil {
		t.Errorf("expected rolling forward to an invalid version number to error")
	}

	// NOW we roll back.
	app.Rollback(1)
	if app.Ledger.Len() != 3 {
		t.Errorf("expected ledger to have 3 releases; got %d", app.Ledger.Len())
	}
	if app.Ledger[2].Build != release1.Build {
		t.Errorf("expected new release build to be v1's build")
	}
	if app.Ledger[2].Version != 3 {
		t.Errorf("expected new release to be v3, got v%d", app.Ledger[2].Version)
	}
}
