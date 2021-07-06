package network

import (
	"net"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type ToNotify struct {
	ToNotifyDocuments chan *utils.PackedDocument
	AddCurrent *utils.PackedDocument
}

type DocumentManager struct {
	Id utils.UUID
	socket net.Conn
	ToNotify ToNotify
}

func (manager *DocumentManager) SetOnMessageReceiveListener(listener MessageReceiveListener)   {
	socket := manager.socket
	for {
		received := make([]byte, 1024)
		_, err := socket.Read(received)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}

		for _, data := range utils.GetPackedDocuments(received) {
			manager.ToNotify.AddCurrent = &data
			listener(manager.ToNotify)
		}
	}
}

func (manager *DocumentManager) SetPeerConnectedListener(listener PeerConnectedListener) {
}

func (manager *DocumentManager) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {
}

func (manager *DocumentManager) BroadcastMessage(message interface{}) {
	packedDocument := message.(utils.PackedDocument)
	bytes := utils.ToBytes(packedDocument)
	sendRequest(manager.socket, bytes)
}

func sendRequest(socket net.Conn, bytes []byte) {
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

// Connect establishes necessary connections and enables
// receiving and sending changes to and from network.
// Applications must set listeners using SetListener
// before calling Start
func (manager *DocumentManager) Connect() {
	for {
		socket, err := net.Dial("tcp", "localhost:8081")
		if err != nil {
			continue
		}
		socket.SetDeadline(time.Now().Add(time.Second))
		manager.socket = socket
		manager.ToNotify = ToNotify{make(chan(*utils.PackedDocument)), nil}
		bytes := utils.ToBytes(utils.PackedDocument{string(manager.Id), "", "", "", "connect"})
		go sendRequest(manager.socket, bytes)
		return 
	}
}

// Disconnect terminates established connections
func (manager *DocumentManager) Disconnect() {
	defer manager.socket.Close()
}

// Kill frees resources and
func (manager *DocumentManager) Kill() {}