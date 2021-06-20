package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Server struct {
	SyncedDocuments []*(crdt.SyncedDocument)
}

var server *Server

func NewServer() *Server {
	if server == nil {
		server = &Server{}
	}
	go server.Listen()
	return server
}

func (server *Server)Listen() {
	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
			// fmt.Println(err)
			return
	}
	defer listener.Close()

	for {
		socket, err := listener.Accept()
		if err != nil {
				// fmt.Println(err)
				return
		}
		socket.SetDeadline(time.Now().Add(time.Second))
		go server.HandleRequest(socket)
	}
}

func insert(data utils.PackedDocument){
	// Send request to all clients
	position := crdt.ToBasicPosition(data.Position)
	value := data.Value

	for _, syncedDoc := range server.SyncedDocuments {
		if strconv.Itoa(syncedDoc.GetSite()) != data.Site {
			syncedDoc.InsertAtPosition(position, value)
		}
	}
}


func delete(data utils.PackedDocument){
	// Send request to all clients
	position := crdt.ToBasicPosition(data.Position)

	for _, syncedDoc := range server.SyncedDocuments {
		if strconv.Itoa(syncedDoc.GetSite()) != data.Site {
			syncedDoc.DeleteAtPosition(position)
		}
	}
}

func (server *Server)ConnectWithClient(doc *crdt.SyncedDocument) {
	server.SyncedDocuments = append(server.SyncedDocuments, doc)
}

func (server *Server) HandleRequest(socket net.Conn) {
	for {
		receivedMessage := make([]byte, 1024)
		_, err := socket.Read(receivedMessage)
		if err != nil {
				fmt.Println(err)
				socket.SetDeadline(time.Now().Add(time.Second))
				continue
		}
		
		if string(receivedMessage) == "Done"{
				continue
		}

		var packedDocument = utils.FromBytes(receivedMessage)
		action := packedDocument.Action

		if action == "Insert" {
			insert(packedDocument)
		} else {
			delete(packedDocument)
		}
		
		sendMessage := "OK"
		socket.Write([]byte(sendMessage))
		
	}
}