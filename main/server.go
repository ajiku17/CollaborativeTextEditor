package main

import (
	"net"
	"sync"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Server struct {
	ConnectedSockets map[net.Conn]string
	Changes []utils.PackedDocument
	lock             *sync.Mutex
}

var server *Server

func NewServer() *Server {
	if server == nil {
		server = &Server{}
		server.ConnectedSockets = make(map[net.Conn]string)
		server.Changes = make([]utils.PackedDocument, 0)
		server.lock = &sync.Mutex{}
	}
	go server.Listen()
	return server
}

func (server *Server) Listen() {
	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		// log.Fatalln(err)
		return
	}
	defer listener.Close()

	for {
		socket, err := listener.Accept()
		if err != nil {
			// log.Fatalln(err)
			return
		}
		socket.SetDeadline(time.Now().Add(time.Second))
		go server.HandleRequest(socket)
	}
}

func (server *Server) sendAll(data utils.PackedDocument) {
	server.lock.Lock()
	for socket, site := range server.ConnectedSockets {
		if site != data.Site {
			socket.Write(utils.ToBytes(data))
		}
	}
	server.lock.Unlock()
}

func (server *Server) HandleRequest(socket net.Conn) {
	for {
		receivedMessage := make([]byte, 1024)
		_, err := socket.Read(receivedMessage)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}

		for _, data := range utils.GetPackedDocuments(receivedMessage) {
			if data.Action == "connect" {
				server.setSocketSite(data.Site, socket)
				go server.syncNewConnection()
			} else {
				server.Changes = append(server.Changes, data)
				go server.sendAll(data)
			}
		}
	}
}

func (server *Server)syncNewConnection() {
	for _, packedDocument := range server.Changes {
		go server.sendAll(packedDocument)
	}
}

func (server *Server) setSocketSite(site string, socket net.Conn) {
	server.ConnectedSockets[socket] = site
}

func (server *Server) IsConnected(site string) bool {
	for _, curr_site := range server.ConnectedSockets {
		if curr_site == site {
			return true
		}
	}
	return false
}