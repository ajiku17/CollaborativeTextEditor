package crdt

type Document interface {
	InsertAtIndex(string, int, int) Position
	DeleteAtIndex(int) Position
	// GetInsertPosition(int, int) Position
	// GetDeletePosition(int) Position
	InsertAtPosition(Position, string)
	DeleteAtPosition(Position)
	ToString() string
	Length() int
}