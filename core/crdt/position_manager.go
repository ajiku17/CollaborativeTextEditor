package crdt

type Position interface{}

type PositionManager interface {
	PositionIsLessThan(Position, Position) bool
	PositionsEqual(Position, Position) bool
	AllocPositionBetween(Position, Position) Position
	GetMaxPosition() Position
	GetMinPosition() Position
}