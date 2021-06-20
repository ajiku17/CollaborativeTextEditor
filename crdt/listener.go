package crdt

type Listener interface {
	AddListener()
	Notify()
}