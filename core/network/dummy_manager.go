package network

type DummyManager struct {

}

func NewDummyManager() Manager {
	manager := new (DummyManager)
	
	return manager
}

func (d *DummyManager) SetOnMessageReceiveListener(listener MessageReceiveListener) {

}

func (d *DummyManager) SetPeerConnectedListener(listener PeerConnectedListener) {

}

func (d *DummyManager) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {

}

func (d *DummyManager) BroadcastMessage(message interface{}) {
}

func (d *DummyManager) Connect() {
	
}

func (d *DummyManager) Disconnect() {

}

func (d *DummyManager) Kill() {

}
