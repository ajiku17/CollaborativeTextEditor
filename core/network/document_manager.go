package network

import (
	"net"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type DocumentManager struct {
	Id utils.UUID
	socket net.Conn
	toSend chan interface{}
	listeners []interface{}   //[MessageReceiveListener, PeerConnectedListener, PeerDisconnectedListener]
}

func NewDocumentManager(id utils.UUID) Manager {
	manager := new (DocumentManager)
	manager.Id = id
	manager.toSend = make(chan(interface{}))
	manager.listeners = make([]interface{}, 3)
	return manager
}

func (manager *DocumentManager) SetOnMessageReceiveListener(listener MessageReceiveListener)   {
	manager.listeners[0] = listener
}

func (manager *DocumentManager) SetPeerConnectedListener(listener PeerConnectedListener) {
}

func (manager *DocumentManager) SetPeerDisconnectedListener(listener PeerDisconnectedListener) {
}

func (manager *DocumentManager) messageRecieved() {
	socket := manager.socket
	for {
		received := make([]byte, 1024)
		_, err := socket.Read(received)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}

		manager.listeners[0].(MessageReceiveListener)(utils.FromBytes(received))
	}
}

func (manager *DocumentManager) BroadcastMessage(message interface{}) {
	manager.toSend <- message
}

func (manager *DocumentManager) broadcastMessages() {
	for {
		message := <- manager.toSend
		bytes := utils.ToBytes(message)
		if bytes != nil {
			go manager.sendRequest(bytes)
		}
	}
}

func (manager *DocumentManager)sendRequest(bytes []byte) {
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
		go manager.messageRecieved()
		go manager.broadcastMessages()
		return 
	}
}

// Disconnect terminates established connections
func (manager *DocumentManager) Disconnect() {
	defer manager.socket.Close()
}

// Kill frees resources and
func (manager *DocumentManager) Kill() {}