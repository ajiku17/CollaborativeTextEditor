package crdt

import (
	"net"
	"strconv"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type BasicDocumentUpdateManager struct {
	url    string
	socket net.Conn
	toNotify *utils.PackedDocument
}

func NewDocumentUpdateManager(serverUrl string) *BasicDocumentUpdateManager {
	for {
		socket, err := net.Dial("tcp", serverUrl)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		socket.SetDeadline(time.Now().Add(time.Second))
		manager := BasicDocumentUpdateManager{serverUrl, socket, nil}
		go manager.AddListener()
		return &manager
	}
}

func (manager *BasicDocumentUpdateManager) Insert(position Position, val string, site int) {
	toBytes := utils.ToBytes(utils.PackedDocument{strconv.Itoa(site), BasicPositionToString(position.(BasicPosition)), val, "Insert"})
	manager.sendRequest(toBytes, "Insert")
}


func (manager *BasicDocumentUpdateManager) Delete(position Position, site int) {
	toBytes := utils.ToBytes(utils.PackedDocument{strconv.Itoa(site), BasicPositionToString(position.(BasicPosition)), "", "Delete"})
	manager.sendRequest(toBytes, "Delete")
}

func (manager *BasicDocumentUpdateManager)sendRequest(data []byte, actionName string){
	socket := manager.socket
	for {
		_, err := socket.Write(data)
		// fmt.Printf("Client send %s", data)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		} else {
			// fmt.Printf("TCP Exit")
			return
		}

		// received := make([]byte, 1024)
		// _, err = socket.Read(received)
		// if err != nil {
		// 	socket.SetDeadline(time.Now().Add(time.Second))
		// 	continue
		// }
		// received = received[:bytes.Index(received, []byte{0})]   //remove trailing zeros
		// if bytes.Compare(received, []byte("OK")) == 0{
		// 		return
		// }
	}
}

func (manager *BasicDocumentUpdateManager)AddListener() {
	socket := manager.socket
	for {
		received := make([]byte, 1024)
		_, err := socket.Read(received)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}
		// fmt.Printf("RECEIVED - %s\n", received)
		var packedDocument = utils.FromBytes(received)
		manager.toNotify = &packedDocument
	}
}


func (manager *BasicDocumentUpdateManager)Notify() *utils.PackedDocument {
	return manager.toNotify
}