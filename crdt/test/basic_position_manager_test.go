package test

import (
	"github.com/ajiku17/CollaborativeTextEditor/crdt"
	"reflect"
	"testing"
)

func TestPositionPrefix(t *testing.T) {
	var position crdt.BasicPosition

	for i := 0; i < 10; i++ {
		position = append(position, crdt.Identifier{Pos: i, Site: i})
	}

	pref := crdt.Prefix(position, 5)
	AssertTrue(t, len(pref) == 5)

	for i := 0; i < 5; i++ {
		AssertTrue(t, pref[i] == position[i].Pos)
	}

	pref = crdt.Prefix(position, 12)
	AssertTrue(t, len(pref) == 12)

	for i := 0; i < 10; i++ {
		AssertTrue(t, pref[i] == position[i].Pos)
	}

	AssertTrue(t, pref[10] == 0)
	AssertTrue(t, pref[11] == 0)

	pref = crdt.Prefix(position, 0)
	AssertTrue(t, len(pref) == 0)
}

func TestPositionToNumber(t *testing.T) {
	var position crdt.BasicPosition

	for i := 0; i < 6; i++ {
		position = append(position, crdt.Identifier{i, i})
	}

	number := crdt.PositionToNumber(position)
	expected := crdt.Number{0, 1, 2, 3, 4, 5}
	AssertTrue(t, reflect.DeepEqual(number, expected))
}

func TestPositionSubtract(t *testing.T) {
	var position1, position2 crdt.BasicPosition

	// #1
	crdt.NumberSetBase(10)
	position1 = append(position1, crdt.Identifier{3, 2})
	position1 = append(position1, crdt.Identifier{8, 2})

	position2 = append(position2, crdt.Identifier{9, 5})
	position2 = append(position2, crdt.Identifier{8, 1})

	AssertTrue(t, crdt.NumberToInt(crdt.PositionSubtract(position2, position1)) == crdt.NumberToInt(crdt.Number{6, 0}))

	// #2
	crdt.NumberSetBase(64)
	position1 = crdt.BasicPosition{}
	position1 = append(position1, crdt.Identifier{3, 2})
	position1 = append(position1, crdt.Identifier{8, 7})

	position2 = crdt.BasicPosition{}
	position2 = append(position2, crdt.Identifier{5, 5})
	position2 = append(position2, crdt.Identifier{6, 1})

	AssertTrue(t, crdt.NumberToInt(crdt.PositionSubtract(position2, position1)) == crdt.NumberToInt(crdt.Number{1, 62}))
}

func TestPositionAdd(t *testing.T) {
	var position1, position2 crdt.BasicPosition

	// #2
	crdt.NumberSetBase(10)
	position1 = append(position1, crdt.Identifier{3, 2})
	position1 = append(position1, crdt.Identifier{9, 3})
	position1 = append(position1, crdt.Identifier{2, 2})

	position2 = append(position2, crdt.Identifier{3, 1})
	position2 = append(position2, crdt.Identifier{9, 5})
	position2 = append(position2, crdt.Identifier{1, 5})

	AssertTrue(t, crdt.NumberToInt(crdt.PositionAdd(position1, position2)) == 783)

	// #2
	crdt.NumberSetBase(64)

	position1 = crdt.BasicPosition{}
	position2 = crdt.BasicPosition{}

	position1 = append(position1, crdt.Identifier{3, 2})
	position1 = append(position1, crdt.Identifier{8, 2})

	position2 = append(position2, crdt.Identifier{9, 5})
	position2 = append(position2, crdt.Identifier{8, 1})

	AssertTrue(t, crdt.NumberToInt(crdt.PositionAdd(position1, position2)) == crdt.NumberToInt(crdt.Number{12, 16}))
}

func TestConstructPosition(t *testing.T) {
	var r crdt.Number
	var prevPos, afterPos crdt.BasicPosition

	prevPos = append(prevPos, crdt.Identifier{3, 2})
	prevPos = append(prevPos, crdt.Identifier{9, 3})
	prevPos = append(prevPos, crdt.Identifier{1, 2})

	afterPos = append(afterPos, crdt.Identifier{3, 1})
	afterPos = append(afterPos, crdt.Identifier{9, 5})
	afterPos = append(afterPos, crdt.Identifier{2, 5})

	r = crdt.Number{3, 9, 1, 4}

	newPos := crdt.ConstructPosition(r, prevPos, afterPos, 8)
	AssertTrue(t, len(newPos) == 4)
	for i := 0; i < len(prevPos); i++ {
		AssertTrue(t, prevPos[i].Pos == newPos[i].Pos && prevPos[i].Site == newPos[i].Site)
	}
	AssertTrue(t, newPos[len(newPos) - 1].Pos == 4 && newPos[len(newPos) - 1].Site == 8)
}

func PositionsEqual(t *testing.T) {
	var position1, position2 crdt.BasicPosition

	// #1
	manager := crdt.NewBasicPositionManager()
	
	crdt.NumberSetBase(10)
	position1 = append(position1, crdt.Identifier{3, 2})
	position1 = append(position1, crdt.Identifier{8, 2})

	position2 = append(position2, crdt.Identifier{9, 5})
	position2 = append(position2, crdt.Identifier{8, 1})

	AssertTrue(t, !manager.PositionsEqual(position1, position2))
	AssertTrue(t, !manager.PositionsEqual(position2, position1))

	// #2
	manager = crdt.NewBasicPositionManager()

	crdt.NumberSetBase(64)
	position1 = crdt.BasicPosition{}
	position1 = append(position1, crdt.Identifier{3, 2})
	position1 = append(position1, crdt.Identifier{8, 7})

	position2 = crdt.BasicPosition{}
	position2 = append(position2, crdt.Identifier{3, 2})
	position2 = append(position2, crdt.Identifier{8, 7})

	AssertTrue(t, manager.PositionsEqual(position1, position2))
	AssertTrue(t, manager.PositionsEqual(position2, position1))

	position2 = append(position2, crdt.Identifier{1, 2})
	AssertTrue(t, manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionsEqual(position1, position2))
	AssertTrue(t, !manager.PositionsEqual(position2, position1))
}

func PositionIsLessThan(t *testing.T) {
	var position1, position2 crdt.BasicPosition
	manager := crdt.NewBasicPositionManager()

	// #1

	crdt.NumberSetBase(10)
	position1 = append(position1, crdt.Identifier{3, 2})
	position1 = append(position1, crdt.Identifier{8, 2})

	position2 = append(position2, crdt.Identifier{9, 5})
	position2 = append(position2, crdt.Identifier{8, 1})

	AssertTrue(t, manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionIsLessThan(position2, position1))

	// #2
	manager = crdt.NewBasicPositionManager()

	crdt.NumberSetBase(64)
	position1 = crdt.BasicPosition{}
	position1 = append(position1, crdt.Identifier{3, 2})
	position1 = append(position1, crdt.Identifier{8, 7})

	position2 = crdt.BasicPosition{}
	position2 = append(position2, crdt.Identifier{3, 2})
	position2 = append(position2, crdt.Identifier{8, 7})

	AssertTrue(t, !manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionIsLessThan(position2, position1))

	position2 = append(position2, crdt.Identifier{1, 2})
	AssertTrue(t, manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionIsLessThan(position2, position1))
}

func PositionAllocation(t *testing.T) {
	var prevPos, afterPos crdt.BasicPosition
	manager := crdt.NewBasicPositionManager()

	// #1
	crdt.NumberSetBase(10)
	prevPos = append(prevPos, crdt.Identifier{7, 2})

	afterPos = append(afterPos, crdt.Identifier{10, 1})

	newP := manager.AllocPositionBetween(prevPos, afterPos, 4)
	newPosition, ok := newP.(crdt.BasicPosition)

	AssertTrue(t, ok)
	AssertTrue(t, len(newPosition) == 1)
	AssertTrue(t, manager.PositionIsLessThan(prevPos, newPosition))
	AssertTrue(t, newPosition[len(newPosition) - 1].Site == 4)
	AssertTrue(t, newPosition[len(newPosition) - 1].Pos > 7)
	AssertTrue(t, newPosition[len(newPosition) - 1].Pos < 10)

	// #2
	manager = crdt.NewBasicPositionManager()

	crdt.NumberSetBase(10)

	prevPos = crdt.BasicPosition{}
	afterPos = crdt.BasicPosition{}

	prevPos = append(prevPos, crdt.Identifier{3, 2})
	prevPos = append(prevPos, crdt.Identifier{9, 3})
	prevPos = append(prevPos, crdt.Identifier{1, 2})

	afterPos = append(afterPos, crdt.Identifier{3, 1})
	afterPos = append(afterPos, crdt.Identifier{9, 5})
	afterPos = append(afterPos, crdt.Identifier{2, 5})

	newP = manager.AllocPositionBetween(prevPos, afterPos, 9)
	newPosition, ok = newP.(crdt.BasicPosition)

	AssertTrue(t, ok)
	AssertTrue(t, len(newPosition) == 4)
	AssertTrue(t, manager.PositionIsLessThan(prevPos, newPosition))
	AssertTrue(t, newPosition[len(newPosition) - 1].Site == 9)
	AssertTrue(t, newPosition[len(newPosition) - 1].Pos < 10)
}