package table

import "sync"

type Table struct {
	table map[string] []string
	mu    sync.Mutex
}

func (t *Table) Register (docId string, addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	peers, ok := t.table[docId]
	if !ok {
		peers = []string{}
	}

	t.table[docId] = append(peers, addr)
}

func (t *Table) Get (docId string) []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	peers, ok := t.table[docId]
	if !ok {
		peers = []string{}
	}

	return peers
}

func New () *Table {
	t := new(Table)

	t.table = make(map[string] []string)

	return t
}