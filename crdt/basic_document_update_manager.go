package crdt

import (
	"bytes"
	"net"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type BasicDocumentUpdateManager struct {
	url    string
	socket net.Conn
}

func NewDocumentUpdateManager(serverUrl string) *BasicDocumentUpdateManager {
	for {
		socket, err := net.Dial("tcp", serverUrl)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		socket.SetDeadline(time.Now().Add(time.Second))
		return &BasicDocumentUpdateManager{serverUrl, socket}
	}
}

func (manager *BasicDocumentUpdateManager) Insert(position Position, val string, site int) {
	toBytes := utils.ToBytes(site, BasicPositionToString(position.(BasicPosition)), val, "Insert")
	manager.sendRequest(toBytes, "Insert")
}


func (manager *BasicDocumentUpdateManager) Delete(position Position, site int) {
	toBytes := utils.ToBytes(site, BasicPositionToString(position.(BasicPosition)), "", "Delete")
	manager.sendRequest(toBytes, "Delete")
}

func (manager *BasicDocumentUpdateManager)sendRequest(data []byte, actionName string){
	socket := manager.socket
	for {
		socket.Write(data)

		received := make([]byte, 1024)
		_, err := socket.Read(received)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}
		received = received[:bytes.Index(received, []byte{0})]   //remove trailing zeros
		if bytes.Compare(received, []byte("OK")) == 0{
				return
		}
	}
}
