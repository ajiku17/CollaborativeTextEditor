package test

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"testing"
)

func CheckPositionOrdering(t *testing.T, manager crdt.PositionManager, prev, mid, after crdt.Position) {
	AssertTrue(t, manager.PositionIsLessThan(prev, mid))
	AssertTrue(t, !manager.PositionIsLessThan(mid, prev))

	AssertTrue(t, !manager.PositionsEqual(mid, after))
	AssertTrue(t, !manager.PositionsEqual(after, mid))

	AssertTrue(t, !manager.PositionsEqual(mid, prev))
	AssertTrue(t, !manager.PositionsEqual(prev, mid))

	AssertTrue(t, manager.PositionIsLessThan(mid, after))
	AssertTrue(t, !manager.PositionIsLessThan(after, mid))
}

func PositionManagerTest(t *testing.T, manager crdt.PositionManager) {
	minPosition := manager.GetMinPosition()
	maxPosition := manager.GetMaxPosition()

	pos1 := manager.AllocPositionBetween(minPosition, maxPosition)
	CheckPositionOrdering(t, manager, minPosition, pos1, maxPosition)

	pos2 := manager.AllocPositionBetween(pos1, maxPosition)
	CheckPositionOrdering(t, manager, pos1, pos2, maxPosition)

	pos3 := manager.AllocPositionBetween(pos1, pos2)
	CheckPositionOrdering(t, manager, pos1, pos3, pos2)

	pos4 := manager.AllocPositionBetween(pos3, pos2)
	CheckPositionOrdering(t, manager, pos3, pos4, pos2)
}


func TestPosition(t *testing.T) {
	implementations :=  []struct {
		newInstance func () crdt.PositionManager
		name string
	} {
		{ func() crdt.PositionManager {
			instance := crdt.NewBasicPositionManager("1")
			return instance
		}, "BasicPositonManager"},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func (t *testing.T) {
			t.Run("TestPositionManager", func (t* testing.T) {
				PositionManagerTest(t, impl.newInstance())
			})
		})
	}
}