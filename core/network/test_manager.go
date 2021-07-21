package network

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"sync"
	"time"
)

type DummyManager struct {
	Id      utils.UUID
	doc     synceddoc.Document
	crdt    crdt.Document
	stopped bool
	mu      sync.Mutex
}

func (d *DummyManager) GetId() utils.UUID {
	return d.Id
}

func (d *DummyManager) Start() {
	fmt.Println("starting manager")
	go d.sync()
}

// random insert or delete ops
func (d *DummyManager) randOp() synceddoc.Op {
	if d.crdt.Length() > 0 {
		n := utils.RandBetween(1, 2)

		if n == 1 {
			val := "a"

			index := utils.RandBetween(0, d.crdt.Length())

			return crdt.OpInsert{
				Pos: d.crdt.InsertAtIndex(val, index),
				Val: val,
			}
		} else {
			index := utils.RandBetween(0, d.crdt.Length() - 1)

			return crdt.OpDelete{
				Pos: d.crdt.DeleteAtIndex(index),
			}
		}
	} else {
		val := "a"

		return crdt.OpInsert{
			Pos: d.crdt.InsertAtIndex(val, 0),
			Val: val,
		}
	}
}

// apply remote op every 2 seconds
func (d *DummyManager) sync () {
	for {
		d.mu.Lock()
		if d.stopped {
			d.mu.Unlock()
			return
		}
		d.mu.Unlock()

		d.doc.ApplyRemoteOp(d.Id, d.randOp(), nil)
		time.Sleep(2 * time.Second)
	}
}

func (d *DummyManager) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.stopped = true
}

func (d *DummyManager) Kill() {
	d.Stop()
}

func NewDummyManager(id string, doc synceddoc.Document) Manager {
	manager := new (DummyManager)
	manager.Id = utils.UUID(id)
	manager.doc = doc
	manager.stopped = false

	manager.crdt = crdt.NewBasicDocument(crdt.NewBasicPositionManager(manager.Id))
	return manager
}