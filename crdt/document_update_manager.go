package crdt

type DocumentUpdateManager interface {
	Insert(position Position, val string, site int)
	Delete(position Position, site int)
}