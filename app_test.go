package api

import (
	"testing"
	"time"
)

func TestCreateApp(t *testing.T) {
	app, _ := NewApp("test")

	// close enough for government work!
	if time.Since(app.Created) > time.Duration(10*time.Millisecond) {
		t.Errorf("Expected the app's created date to be within 1ms, got '%v'", time.Since(app.Created))
	}

	if time.Since(app.Updated) > time.Duration(10*time.Millisecond) {
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
	app, _ := NewApp("")

	if app.ID == "" {
		t.Error("expected app ID to be generated, got empty string")
	}
}

func TestAppRelease(t *testing.T) {
	app, _ := NewApp("")
	release := app.NewRelease(&Build{}, &Config{})
	if release == nil {
		t.Errorf("expected app to create a new release")
	}
	// creating an app creates an initial release
	if release.Version != 2 {
		t.Errorf("expected version to be 2; got %d", release.Version)
	}
	if app.Ledger.Len() != 2 {
		t.Errorf("expected release to be appended to the ledger; got %d", app.Ledger.Len())
	}
	app.NewRelease(&Build{}, &Config{})
	if app.Ledger.Len() != 3 {
		t.Errorf("expected ledger to have 3 releases; got %d", app.Ledger.Len())
	}
	if !app.Ledger.Less(1, 2) {
		t.Errorf("expected v2 to be less than v3. v1 = %d; v2 = %d", app.Ledger[1].Version, app.Ledger[2].Version)
	}
	app.Ledger.Swap(1, 2)
	if app.Ledger[0] != release {
		t.Errorf("expected v2 to be first; got %s", app.Ledger[0])
	}
}

func TestAppRollback(t *testing.T) {
	app, _ := NewApp("")
	release2 := app.NewRelease(&Build{}, &Config{})
	app.NewRelease(&Build{}, &Config{})

	// first, check that we cannot roll back to an invalid version
	if err := app.Rollback(0); err == nil {
		t.Errorf("expected rolling back to an invalid version number to error")
	}

	// now check that we cannot roll forward to a release that does not exist yet
	if err := app.Rollback(100); err == nil {
		t.Errorf("expected rolling forward to an invalid version number to error")
	}

	// NOW we roll back.
	app.Rollback(2)
	if app.Ledger.Len() != 4 {
		t.Errorf("expected ledger to have 4 releases; got %d", app.Ledger.Len())
	}
	if app.Ledger[len(app.Ledger)-1].Build != release2.Build {
		t.Errorf("expected new release build to be v1's build")
	}
	if app.Ledger[len(app.Ledger)-1].Version != 4 {
		t.Errorf("expected new release to be v4, got v%d", app.Ledger[2].Version)
	}
}
