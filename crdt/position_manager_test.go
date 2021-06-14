package crdt

import (
	"testing"
)

func PositionsEqual(t *testing.T,  manager PositionManager) {
	var position1, position2 BasicPosition

	// #1
	manager.PositionManagerInit()
	
	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, !manager.PositionsEqual(position1, position2))
	AssertTrue(t, !manager.PositionsEqual(position2, position1))

	// #2
	manager.PositionManagerInit()

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

func PositionIsLessThan(t *testing.T, manager PositionManager) {
	var position1, position2 BasicPosition
	manager.PositionManagerInit()

	// #1
	manager.PositionManagerInit()

	NumberSetBase(10)
	position1 = append(position1, Identifier{3, 2})
	position1 = append(position1, Identifier{8, 2})

	position2 = append(position2, Identifier{9, 5})
	position2 = append(position2, Identifier{8, 1})

	AssertTrue(t, manager.PositionIsLessThan(position1, position2))
	AssertTrue(t, !manager.PositionIsLessThan(position2, position1))

	// #2
	manager.PositionManagerInit()

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

func PositionAllocation(t *testing.T, manager PositionManager) {
	var prevPos, afterPos BasicPosition
	manager.PositionManagerInit()

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
	manager.PositionManagerInit()

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

func TestPosition(t *testing.T) {
	implementations :=  []struct {
		newInstance func () PositionManager
		name string
	} {
		{ func() PositionManager {
			return new(BasicPositionManager)
		}, "BasicPositonManager"},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func (t *testing.T) {
			t.Run("TestPositionAllocation", func (t* testing.T) {
				PositionAllocation(t, impl.newInstance())
			})
			t.Run("TestPositionsEqual", func (t* testing.T) {
				PositionsEqual(t, impl.newInstance())
			})
			t.Run("TestPositionIsLessThan", func (t* testing.T) {
				PositionIsLessThan(t, impl.newInstance())
			})
		})
	}
}