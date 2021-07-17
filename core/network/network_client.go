package network

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"net"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type NetworkClient struct {
	id utils.UUID
	socket net.Conn
	document *synceddoc.Document
	alive bool
}

func NewDocumentManager(id utils.UUID, synccedDoc *synceddoc.Document) Manager {
	manager := new (NetworkClient)
	manager.id = id
	manager.document = synccedDoc
	manager.alive = false
	manager.connect()
	return manager
}

func (manager *NetworkClient) GetId() utils.UUID {
	return manager.id
}

func (manager *NetworkClient) Start() {
	manager.alive = true
	operationRequest := synceddoc.OperationRequest{manager.id, synceddoc.ConnectRequest{manager.id}}
	go manager.sendRequest(utils.ToBytes(operationRequest))
	go manager.messageReceived()
	go manager.broadcastMessages()
}

func (manager *NetworkClient) Stop() {
	manager.alive = false
}

// Kill frees resources and
func (manager *NetworkClient) Kill() {
	defer manager.socket.Close()
}


// Connect establishes necessary connections and enables
// receiving and sending changes to and from network.
// Applications must set listeners using SetListener
// before calling Start
func (manager *NetworkClient) connect() {
	for {
		socket, err := net.Dial("tcp", "localhost:8081")
		if err != nil {
			continue
		}
		socket.SetDeadline(time.Now().Add(time.Second))
		manager.socket = socket
		return
	}
}

func (manager *NetworkClient) messageReceived() {
	socket := manager.socket
	for manager.alive {
		received := make([]byte, 1024)
		_, err := socket.Read(received)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}
		data := utils.FromBytes(received).(synceddoc.OperationRequest)
		op := data.Operation
		(*manager.document).ApplyRemoteOp(data.Id, op, nil)
	}
}

func (manager *NetworkClient) broadcastMessages() {
	for manager.alive {
		message := (*manager.document).NextUnbroadcastedChange()
		if message == nil {
			continue
		}
		bytes := utils.ToBytes(message)
		if bytes != nil {
			go manager.sendRequest(bytes)
		}
	}
}

func (manager *NetworkClient)sendRequest(bytes []byte) {
	socket := manager.socket
	for {
		_, err := socket.Write(bytes)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		} else {
			return
		}
	}
}