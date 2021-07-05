package crdt

type Document interface {
	InsertAtIndex(string, int, int) Position
	DeleteAtIndex(int)
	InsertAtPosition(Position, string)
	DeleteAtPosition(Position)
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	ToString() string
	Length() int
}
