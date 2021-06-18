package crdt

type Position interface{}

type PositionManager interface {
	PositionManagerInit()	
	PositionIsLessThan(Position, Position) bool
	PositionsEqual(Position, Position) bool
	AllocPositionBetween(Position, Position, int) Position
	GetMaxPosition() Position
	GetMinPosition() Position
}