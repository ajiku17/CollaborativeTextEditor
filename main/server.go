package main

import (
	"fmt"
	"net"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Server struct {
	SyncedDocuments []*(crdt.SyncedDocument)
	connectedSockets map[net.Conn]string
}

var server *Server

func NewServer() *Server {
	if server == nil {
		server = &Server{}
	}
	server.connectedSockets = make(map[net.Conn]string)
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
		server.connectedSockets[socket] = ""
		go server.HandleRequest(socket)
	}
}

func sendAll(data utils.PackedDocument){
	// Send request to all clients
	for socket, site := range server.connectedSockets {
		if site != data.Site {
			fmt.Printf("Server Send %s\n", utils.ToBytes(data))
			socket.Write(utils.ToBytes(data))
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
				socket.SetDeadline(time.Now().Add(time.Second))
				continue
		}
		
		if string(receivedMessage) == "Done"{
				continue
		}

		fmt.Printf("Server received %s\n", receivedMessage)

		var packedDocument = utils.FromBytes(receivedMessage)

		server.setSocketSite(packedDocument.Site, socket)
		
		sendAll(packedDocument)
		
	}
}

func (server *Server) setSocketSite(site string, socket net.Conn) {
	server.connectedSockets[socket] = site
}