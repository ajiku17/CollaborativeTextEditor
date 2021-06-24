package crdt

type DocumentID uint64

type Document interface {
	DocumentID() DocumentID
	InsertAtIndex(string, int, int) Position
	DeleteAtIndex(int)
	InsertAtPosition(Position, string)
	DeleteAtPosition(Position)
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	ToString() string
	Length() int
}
