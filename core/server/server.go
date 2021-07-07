package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Server struct {
	ConnectedSockets map[utils.UUID]net.Conn
	Changes []interface{}
	lock             *sync.Mutex
}

var server *Server

func NewServer() *Server {
	if server == nil {
		server = &Server{}
		server.ConnectedSockets = make(map[utils.UUID]net.Conn)
		server.Changes = make([]interface{}, 0)
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

func (server *Server) sendAll(currId utils.UUID, data interface{}) {
	server.lock.Lock()
	defer server.lock.Unlock()
	for id, socket := range server.ConnectedSockets {
		if id != currId {
			socket.Write(utils.ToBytes(data))
		}
	}
}

func (server *Server) sendNotify(id utils.UUID, data interface{}) {
	server.lock.Lock()
	defer server.lock.Unlock()
	server.ConnectedSockets[id].Write(utils.ToBytes(data))
}

func (server *Server) HandleRequest(socket net.Conn) {
	for {
		receivedMessage := make([]byte, 1024)
		_, err := socket.Read(receivedMessage)
		if err != nil {
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}

		data := utils.FromBytes(receivedMessage)


		switch data.(type) {
			case synceddoc.ConnectRequest:
				server.setSocketId(data.(synceddoc.ConnectRequest).Id, socket)
				go server.syncNewConnection(data.(synceddoc.ConnectRequest).Id)
			case synceddoc.ChangeCRDTInsert:
				server.Changes = append(server.Changes, data)
				go server.sendAll(data.(synceddoc.ChangeCRDTInsert).ManagerId, data)
			case synceddoc.ChangeCRDTDelete:
				server.Changes = append(server.Changes, data)
				go server.sendAll(data.(synceddoc.ChangeCRDTDelete).ManagerId, data)
			default:
				fmt.Println("Different type")
		}
	}
}

func (server *Server)syncNewConnection(id utils.UUID) {
	for _, bytes := range server.Changes {
		go server.sendNotify(id, bytes)
	}
}

func (server *Server) setSocketId(id utils.UUID, socket net.Conn) {
	server.ConnectedSockets[id] = socket
}


func (server *Server) IsConnected(id utils.UUID) bool {
	for curr_id, _ := range server.ConnectedSockets {
		if curr_id == id {
			return true
		}
	}
	return false
}