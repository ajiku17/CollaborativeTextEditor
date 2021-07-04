package crdt

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Action int

const (
    Insert Action = 0
    Delete Action = 1
)

type BasicDocumentUpdateManager struct {
	url    string
	socket net.Conn
	toNotify chan *utils.PackedDocument
}

func NewDocumentUpdateManager(serverUrl string) *BasicDocumentUpdateManager {
	for {
		socket, err := net.Dial("tcp", serverUrl)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		socket.SetDeadline(time.Now().Add(time.Second))
		toNotify := make(chan(*utils.PackedDocument))
		manager := BasicDocumentUpdateManager{serverUrl, socket, toNotify}
		go manager.AddListener()
		return &manager
	}
}

func (manager *BasicDocumentUpdateManager) ConnectWithServer(site int) {
	toBytes := utils.ToBytes(utils.PackedDocument{strconv.Itoa(site), "", "", "Connect"})
	manager.sendRequest(toBytes)
}

func (manager *BasicDocumentUpdateManager) Insert(position Position, val string, site int) {
	toBytes := utils.ToBytes(utils.PackedDocument{strconv.Itoa(site), BasicPositionToString(position.(BasicPosition)), val, "Insert"})
	manager.sendRequest(toBytes)
}


func (manager *BasicDocumentUpdateManager) Delete(position Position, site int) {
	toBytes := utils.ToBytes(utils.PackedDocument{strconv.Itoa(site), BasicPositionToString(position.(BasicPosition)), "", "Delete"})
	manager.sendRequest(toBytes)
}

func (manager *BasicDocumentUpdateManager)sendRequest(data []byte){
	socket := manager.socket
	for {
		_, err := socket.Write(data)
		// fmt.Printf("Client send %s\n", data)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		} else {
			// fmt.Printf("TCP Exit")
			return
		}
	}
}

func (manager *BasicDocumentUpdateManager)AddListener() {
	socket := manager.socket
	fmt.Printf("Listening to %s\n", socket)
	for {
		received := make([]byte, 1024)
		_, err := socket.Read(received)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}

		fmt.Printf("Listener Received %s\n", received)

		for _, data := range utils.GetPackedDocuments(received) {
			manager.toNotify <- &data
		}
	}
}


func (manager *BasicDocumentUpdateManager)Notify() *utils.PackedDocument {
	return <- manager.toNotify
}