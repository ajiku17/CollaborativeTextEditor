package network

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/tracker"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"sync"
	"time"
)

type DummyManager struct {
	Id          utils.UUID
	lastOpIndex int
	doc         synceddoc.Document
	crdt        crdt.Document
	stopped     bool
	mu          sync.Mutex
}

func (d *DummyManager) ConnectSignals(peerConnectedListener PeerConnectedListener, peerDisconnectedListener PeerDisconnectedListener) {

}

func (d *DummyManager) OnPeerConnect(peerConnectedListener PeerConnectedListener) {

}

func (d *DummyManager) OnPeerDisconnect(peerDisconnectedListener PeerDisconnectedListener) {

}

func (d *DummyManager) GetId() utils.UUID {
	return d.Id
}

func (d *DummyManager) Start() {
	fmt.Println("starting manager")

	client := tracker.NewClient("http://127.0.0.1:9090")

	fmt.Println("your id", d.Id)
	client.Register("dokydok", string(d.Id))
	get, _ := client.Get("dokydok")
	fmt.Println("get peers", get)
	go d.sync()
}

// random insert or delete ops
func (d *DummyManager) randOp() synceddoc.Op {
	var cmd interface{}
	if d.crdt.Length() > 0 {
		n := utils.RandBetween(1, 2)

		if n == 1 {
			val := "a"

			index := utils.RandBetween(0, d.crdt.Length())

			cmd = crdt.OpInsert{
				Pos: d.crdt.InsertAtIndex(val, index),
				Val: val,
			}
		} else {
			index := utils.RandBetween(0, d.crdt.Length() - 1)

			cmd = crdt.OpDelete{
				Pos: d.crdt.DeleteAtIndex(index),
			}
		}
	} else {
		val := "a"

		cmd = crdt.OpInsert{
			Pos: d.crdt.InsertAtIndex(val, 0),
			Val: val,
		}
	}

	d.lastOpIndex++
	return synceddoc.Op{
		PeerId:      d.Id,
		PeerOpIndex: d.lastOpIndex,
		Cmd:         cmd,
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

		d.doc.ApplyRemoteOp(d.randOp(), nil)
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