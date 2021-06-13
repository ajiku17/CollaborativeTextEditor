package crdt

type Document interface {
	DocumentInit(PositionManager)
	InsertAtIndex(string, int, int) interface{}
	DeleteAtIndex(int)
	InsertAtPosition(interface{}, string)
	DeleteAtPosition(interface{})
	ToString() string
	Length() int
}
