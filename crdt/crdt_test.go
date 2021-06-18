package crdt

import (
	"testing"
)

func LongsEqual(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Helper()
		t.Errorf("expected %v actual %v!", expected, actual)
	}
}

func AssertTrue(t *testing.T, condition bool) {
	if !condition {
		t.Helper()
		t.Errorf("assertion failed")
	}
}