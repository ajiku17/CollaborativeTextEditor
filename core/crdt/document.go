package crdt

type Document interface {
	InsertAtIndex(string, int) Position
	DeleteAtIndex(int) Position
	InsertAtPosition(Position, string) int
	DeleteAtPosition(Position) int
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	ToString() string
	Length() int
	GetNextHistoryData(index int) interface{}
	GetHistory() []interface{}
}
