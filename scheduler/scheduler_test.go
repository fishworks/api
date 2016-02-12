package scheduler

import (
	"testing"
)

func TestNew(t *testing.T) {
	if _, err := New("foo"); err == nil {
		t.Errorf("expected scheduler to error with unknown scheduler module")
	}
}
