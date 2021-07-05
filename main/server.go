package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Server struct {
	ConnectedSockets map[net.Conn]string
	lock             *sync.Mutex
}

var server *Server

func NewServer() *Server {
	if server == nil {
		server = &Server{}
		server.ConnectedSockets = make(map[net.Conn]string)
		server.lock = &sync.Mutex{}
	}
	go server.Listen()
	return server
}

func (server *Server) Listen() {
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

func (server *Server) sendAll(data utils.PackedDocument) {
	// Send request to all clients
	server.lock.Lock()
	for socket, site := range server.ConnectedSockets {
		// fmt.Printf("trying to send(range %d), curr site - %d, foreach site - %d, value - %s\n", len(server.ConnectedSockets), data.Site, site, data.Value)
		if site != data.Site {
			fmt.Printf("Server Send %s into %s\n", utils.ToBytes(data), site)
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

		// fmt.Printf("Server received %s\n", receivedMessage)
		for _, data := range utils.GetPackedDocuments(receivedMessage) {
			// fmt.Printf("Server received as packedDocument %s\n", data)
			if data.Action == "connect" {
				server.setSocketSite(data.Site, socket)
			} else {
				go server.sendAll(data)
			}
		}
	}
}

func (server *Server) setSocketSite(site string, socket net.Conn) {
	server.ConnectedSockets[socket] = site
	// fmt.Printf("Connected clients  - %s ; size - %d\n", server.ConnectedSockets, len(server.ConnectedSockets))
}

func (server *Server) IsConnected(site string) bool {
	for _, curr_site := range server.ConnectedSockets {
		if curr_site == site {
			return true
		}
	}
	return false
}