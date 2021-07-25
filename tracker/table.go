package tracker

import "sync"

type Table struct {
	table map[string] map[string] struct {}
	mu    sync.Mutex
}

func (t *Table) Register (docId string, addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, ok := t.table[docId]
	if !ok {
		t.table[docId] = make(map[string] struct{})
	}

	t.table[docId][addr] = struct{}{}
}
func (t *Table) RegisterAndGet (docId string, addr string) []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, ok := t.table[docId]
	if !ok {
		t.table[docId] = make(map[string] struct{})
	}

	t.table[docId][addr] = struct{}{}

	peers := []string {}
	for p := range t.table[docId] {
		peers = append(peers, p)
	}

	return peers
}

func (t *Table) Get (docId string) []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	peers := []string {}
	for p := range t.table[docId] {
		peers = append(peers, p)
	}

	return peers
}

func NewTable () *Table {
	t := new(Table)

	t.table = make(map[string] map[string] struct{})

	return t
}