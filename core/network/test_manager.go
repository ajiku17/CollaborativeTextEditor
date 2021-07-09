package network

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type DummyManager struct {
	DocumentManager
}

func NewDummyManager(id utils.UUID) Manager {
	manager := new (DummyManager)
	manager.Id = id
	manager.toSend = make(chan(interface{}))
	manager.listeners = make([]interface{}, 3)
	return manager
}

func (d *DummyManager) SetOnMessageReceiveListener(listener MessageReceiveListener) {

}

func (d *DummyManager) SetPeerConnectedListener(listener PeerConnectedListener) {

}

func (d *DummyManager) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {

}

// func (d *DummyManager) BroadcastMessage(message interface{}) {
// }

// func (d *DummyManager) Connect() {
// 	gob.Register(ConnectRequest{})
// 	d.DocumentManager.Connect()
// 	d.DocumentManager.BroadcastMessage(ConnectRequest{d.GetId()})
// }

func (d *DummyManager) Disconnect() {

}

func (d *DummyManager) Kill() {

}

func (d *DummyManager) GetId() utils.UUID {
	return d.Id
}