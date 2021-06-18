package crdt

import (
	"reflect"
	"testing"
)

func TestPositionPrefix(t *testing.T) {
	var position BasicPosition

	for i := 0; i < 10; i++ {
		position = append(position, Identifier{pos: i, site: i})
	}

	pref := prefix(position, 5)
	AssertTrue(t, len(pref) == 5)

	for i := 0; i < 5; i++ {
		AssertTrue(t, pref[i] == position[i].pos)
	}

	pref = prefix(position, 12)
	AssertTrue(t, len(pref) == 12)

	for i := 0; i < 10; i++ {
		AssertTrue(t, pref[i] == position[i].pos)
	}

	AssertTrue(t, pref[10] == 0)
	AssertTrue(t, pref[11] == 0)

	pref = prefix(position, 0)
	AssertTrue(t, len(pref) == 0)
}

func TestPositionToNumber(t *testing.T) {
	var position BasicPosition

	for i := 0; i < 6; i++ {
		position = append(position, Identifier{i, i})
	}

	number := PositionToNumber(position)
	expected := Number{0, 1, 2, 3, 4, 5}
	AssertTrue(t, reflect.DeepEqual(number, expected))
}

func TestPositionSubtract(t *testing.T) {
	var position1, position2 BasicPosition

	// #1
	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, NumberToInt(PositionSubtract(position2, position1)) == NumberToInt(Number{6, 0}))

	// #2
	NumberSetBase(64)
	position1 = BasicPosition{}
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 7})

	position2 = BasicPosition{}
	position2 = append(position2, Identifier{5, 5})
	position2 = append(position2, Identifier{6, 1})

	AssertTrue(t, NumberToInt(PositionSubtract(position2, position1)) == NumberToInt(Number{1, 62}))
}

func TestPositionAdd(t *testing.T) {
	var position1, position2 BasicPosition

	// #2
	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{9, 3})
	position1 = append(position1, Identifier{2, 2})

	position2 = append(position2, Identifier{3, 1})
	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{1, 5})

	AssertTrue(t, NumberToInt(PositionAdd(position1, position2)) == 783)

	// #2
	NumberSetBase(64)

	position1 = BasicPosition{}
	position2 = BasicPosition{}

	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, NumberToInt(PositionAdd(position1, position2)) == NumberToInt(Number{12, 16}))
}

func TestConstructPosition(t *testing.T) {
	var r Number
	var prevPos, afterPos BasicPosition

	prevPos = append(prevPos, Identifier{3, 2})
	prevPos = append(prevPos, Identifier{9, 3})
	prevPos = append(prevPos, Identifier{1, 2})

	afterPos = append(afterPos, Identifier{3, 1})
	afterPos = append(afterPos, Identifier{9, 5})
	afterPos = append(afterPos, Identifier{2, 5})

	r = Number{3, 9, 1, 4}

	newPos := constructPosition(r, prevPos, afterPos, 8)
	AssertTrue(t, len(newPos) == 4)
	for i := 0; i < len(prevPos); i++ {
		AssertTrue(t, prevPos[i].pos == newPos[i].pos && prevPos[i].site == newPos[i].site)
	}
	AssertTrue(t, newPos[len(newPos) - 1].pos == 4 && newPos[len(newPos) - 1].site == 8)
}

func PositionsEqual(t *testing.T) {
	var position1, position2 BasicPosition

	// #1
	manager := NewBasicPositionManager()
	
	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, !manager.PositionsEqual(position1, position2))
	AssertTrue(t, !manager.PositionsEqual(position2, position1))

	// #2
	manager = NewBasicPositionManager()

	NumberSetBase(64)
	position1 = BasicPosition{}
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 7})

	position2 = BasicPosition{}
	position2 = append(position2, Identifier{3, 2})
	position2 = append(position2, Identifier{8, 7})

	AssertTrue(t, manager.PositionsEqual(position1, position2))
	AssertTrue(t, manager.PositionsEqual(position2, position1))

	position2 = append(position2, Identifier{1, 2})
	AssertTrue(t, manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionsEqual(position1, position2))
	AssertTrue(t, !manager.PositionsEqual(position2, position1))
}

func PositionIsLessThan(t *testing.T) {
	var position1, position2 BasicPosition
	manager := NewBasicPositionManager()

	// #1

	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionIsLessThan(position2, position1))

	// #2
	manager = NewBasicPositionManager()

	NumberSetBase(64)
	position1 = BasicPosition{}
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 7})

	position2 = BasicPosition{}
	position2 = append(position2, Identifier{3, 2})
	position2 = append(position2, Identifier{8, 7})

	AssertTrue(t, !manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionIsLessThan(position2, position1))

	position2 = append(position2, Identifier{1, 2})
	AssertTrue(t, manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionIsLessThan(position2, position1))
}

func PositionAllocation(t *testing.T) {
	var prevPos, afterPos BasicPosition
	manager := NewBasicPositionManager()

	// #1
	NumberSetBase(10)
	prevPos = append(prevPos, Identifier{7, 2})

	afterPos = append(afterPos, Identifier{10, 1})

	newP := manager.AllocPositionBetween(prevPos, afterPos, 4)
	newPosition, ok := newP.(BasicPosition)

	AssertTrue(t, ok)
	AssertTrue(t, len(newPosition) == 1)
	AssertTrue(t, manager.PositionIsLessThan(prevPos, newPosition))
	AssertTrue(t, newPosition[len(newPosition) - 1].site == 4)
	AssertTrue(t, newPosition[len(newPosition) - 1].pos > 7)
	AssertTrue(t, newPosition[len(newPosition) - 1].pos < 10)

	// #2
	manager = NewBasicPositionManager()

	NumberSetBase(10)

	prevPos = BasicPosition{}
	afterPos = BasicPosition{}

	prevPos = append(prevPos, Identifier{3, 2})
	prevPos = append(prevPos, Identifier{9, 3})
	prevPos = append(prevPos, Identifier{1, 2})

	afterPos = append(afterPos, Identifier{3, 1})
	afterPos = append(afterPos, Identifier{9, 5})
	afterPos = append(afterPos, Identifier{2, 5})

	newP = manager.AllocPositionBetween(prevPos, afterPos, 9)
	newPosition, ok = newP.(BasicPosition)

	AssertTrue(t, ok)
	AssertTrue(t, len(newPosition) == 4)
	AssertTrue(t, manager.PositionIsLessThan(prevPos, newPosition))
	AssertTrue(t, newPosition[len(newPosition) - 1].site == 9)
	AssertTrue(t, newPosition[len(newPosition) - 1].pos < 10)
}