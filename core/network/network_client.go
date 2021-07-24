package network

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/tracker/client"
	"net"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type NetworkClient struct {
	id utils.UUID
	socket net.Conn
	document *synceddoc.Document
	trackerClient *client.Client
	alive bool
}

func NewDocumentManager(id utils.UUID, synccedDoc *synceddoc.Document) Manager {
	manager := new (NetworkClient)
	manager.id = id
	manager.document = synccedDoc
	manager.alive = false
	manager.trackerClient = client.New("localhost:8080")
	manager.connect()
	return manager
}

func (manager *NetworkClient) GetId() utils.UUID {
	return manager.id
}

func (manager *NetworkClient) Start() {
	var operationRequest synceddoc.OperationRequest
	manager.alive = true
	if(len((*manager.document).GetLogs()) == 0) {
		operationRequest = synceddoc.OperationRequest{manager.id, synceddoc.ConnectRequest{make(map[utils.UUID]int)}}
	} else {
		operationRequest = synceddoc.OperationRequest{manager.id, synceddoc.ConnectRequest{(*manager.document).GetLogs()[0].LogState}}
	}
	manager.sendRequest(utils.ToBytes(operationRequest))

	//peers := manager.trackerClient.Get(string((*manager.document).Get    ID()))
	//if len(peers) == 0 {
	//	manager.trackerClient.Register(string((*manager.document).GetID()), string(manager.GetId()))
	//}

	go manager.messageReceived()
	go manager.broadcastMessages()
}

func (manager *NetworkClient) Stop() {
	manager.alive = false
}

// Kill frees resources and
func (manager *NetworkClient) Kill() {
	manager.alive = false
	operationRequest := synceddoc.OperationRequest{manager.id, synceddoc.DisconnectRequest{manager.id}}
	manager.sendRequest(utils.ToBytes(operationRequest))
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
		nextChange := (*manager.document).NextUnbroadcastedChange()
		if nextChange == nil {
			continue
		}
		bytes := utils.ToBytes(synceddoc.OperationRequest{manager.id, nextChange})
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
			//fmt.Println(err)
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		} else {
			return
		}
	}
}

func (manager *NetworkClient) IsAlive() bool{
	return manager.alive
}