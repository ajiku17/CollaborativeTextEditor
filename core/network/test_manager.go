package network


//type TestManager struct {
//	id utils.UUID
//	document *synceddoc.Document
//	alive bool
//	connected bool
//}
//
//func NewTestManager(id utils.UUID, synccedDoc *synceddoc.Document) Manager {
//	manager := new (TestManager)
//	manager.id = id
//	manager.document = synccedDoc
//	return manager
//}
//
//func (d TestManager) GetId() utils.UUID {
//	return d.id
//}
//
//func (d TestManager) Start() {
//	d.alive = true
//	d.connected = true
//	//go d.messageReceived()
//	//go d.broadcastMessages()
//}
//
//func (d TestManager) Stop() {
//	d.alive = false
//	d.connected = true
//}
//
//func (d TestManager) Kill() {
//	d.alive = false
//	d.connected = false
//}

func (d NetworkClient) IsAlive() bool {
	return d.alive
}

//
//func (d TestManager) messageReceived() {
//	//socket := manager.socket
//	//for manager.alive {
//	//	received := make([]byte, 1024)
//	//	_, err := socket.Read(received)
//	//	if err != nil {
//	//		socket.SetDeadline(time.Now().Add(time.Second))
//	//		continue
//	//	}
//	//	data := utils.FromBytes(received).(synceddoc.OperationRequest)
//	//	op := data.Operation
//	//	(*manager.document).ApplyRemoteOp(data.Id, op, nil)
//	//}
//}
//
//func (d TestManager) broadcastMessages() {
//	//for manager.alive {
//	//	nextChange := (*manager.document).NextUnbroadcastedChange()
//	//	if nextChange == nil {
//	//		continue
//	//	}
//	//	fmt.Println("sending ", nextChange)
//	//	bytes := utils.ToBytes(synceddoc.OperationRequest{manager.id, nextChange})
//	//	if bytes != nil {
//	//		go manager.sendRequest(bytes)
//	//	}
//	//}
//}