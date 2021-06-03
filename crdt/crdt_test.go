package crdt

import (
	"testing"
)

func TestPositionPrefix(t *testing.T) {
	var position Position
	zero := Identifier{0, 0}

	for i := 0; i < 10; i++ {
		position = append(position, Identifier{pos: i, site: i})
	}

	pref := prefix(position, 5)
	if len(pref) != 5 {
		t.Error("invalid size")
	}

	for i := 0; i < 5; i++ {
		if pref[i] != position[i] {
			t.Errorf("expected %v got %v!", position[i], pref[i])
		}
	}

	pref = prefix(position, 12)
	if len(pref) != 12 {
		t.Error("invalid size")
	}

	for i := 0; i < 10; i++ {
		if pref[i] != position[i] {
			t.Errorf("expected %v got %v!", position[i], pref[i])
		}
	}

	if pref[10] != zero {
		t.Errorf("expected %v got %v!", position[10], zero)
	}

	if pref[11] != zero {
		t.Errorf("expected %v got %v!", position[11], zero)
	}

	pref = prefix(position, 0)
	if len(pref) != 0 {
		t.Error("invalid size")
	}
}
