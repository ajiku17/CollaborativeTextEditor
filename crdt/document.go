package crdt

type Document interface {
	DocumentInit(PositionManager)
	InsertAtIndex(string, int, int) Position
	DeleteAtIndex(int)
	InsertAtPosition(Position, string)
	DeleteAtPosition(Position)
	ToString() string
	Length() int
}
