package server

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"net"
	"sync"
	"time"

	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Server struct {
	ConnectedSockets map[utils.UUID]net.Conn
	Changes          []interface{}
	GlobalLogCountState map[utils.UUID] int
	GlobalLogState map[utils.UUID] []synceddoc.OperationRequest
	lock             *sync.Mutex
}

var server *Server

func NewServer() *Server {
	if server == nil {
		server = &Server{}
		server.ConnectedSockets = make(map[utils.UUID]net.Conn)
		server.Changes = make([]interface{}, 0)
		server.GlobalLogCountState = make(map[utils.UUID]int)
		server.GlobalLogState = make(map[utils.UUID] []synceddoc.OperationRequest)
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

func (server *Server) sendAll(data synceddoc.OperationRequest) {
	server.lock.Lock()
	defer server.lock.Unlock()
	for currId, socket := range server.ConnectedSockets {
		if data.Id != currId {
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
		//fmt.Println("SErver reading")
		if err != nil {
			//fmt.Println(err)
			socket.SetDeadline(time.Now().Add(time.Second))
			continue
		}

		data := utils.FromBytes(receivedMessage).(synceddoc.OperationRequest)

		operation := data.Operation
		//fmt.Println("SErver got ", operation)

		switch operation.(type) {
			case synceddoc.ConnectRequest:
				server.setSocketId(data.Id, socket)
				go server.syncNewConnection(data.Id, operation.(synceddoc.ConnectRequest))
			case synceddoc.DisconnectRequest:
				server.removeSocketId(operation.(synceddoc.DisconnectRequest).Id)
			case crdt.OpInsert:
				server.Changes = append(server.Changes, data)
				server.incrementLogSequence(data.Id, data)
				go server.sendAll(data)
			case crdt.OpDelete:
				server.Changes = append(server.Changes, data)
				server.incrementLogSequence(data.Id, data)
				go server.sendAll(data)
			default:
				fmt.Println("Different type")
		}
	}
}

func (server *Server)syncNewConnection(id utils.UUID, operation synceddoc.ConnectRequest) {
	fmt.Println("Sync")
	peerLog := operation.PrevLog
	for peerId, globalCount := range server.GlobalLogCountState {
		if count, ok := peerLog[peerId]; ok {
			if globalCount != count {
				fmt.Println("Sending last ", globalCount - count)
				go server.sendLastNOperations(id, peerId, globalCount - count)
			}
		} else {
			fmt.Println("Sending last ", globalCount, "(all)")
			go server.sendLastNOperations(id, peerId, globalCount)
		}
	}
}

func (server *Server) sendLastNOperations(id utils.UUID, operationOwnerId utils.UUID, n int) {
	globalLog := server.GlobalLogState
	for _, bytes := range globalLog[operationOwnerId][len(globalLog[operationOwnerId]) - n:] {
		go server.sendNotify(id, bytes)
	}
}

func (server *Server) setSocketId(id utils.UUID, socket net.Conn) {
	server.ConnectedSockets[id] = socket
}


func (server *Server) removeSocketId(id utils.UUID) {
	delete(server.ConnectedSockets, id)
}

func (server *Server) incrementLogSequence(peerId utils.UUID, data synceddoc.OperationRequest) {
	fmt.Println("Increment ", peerId)
	if _, ok := server.GlobalLogState[peerId]; ok {
		server.GlobalLogState[peerId] = append(server.GlobalLogState[peerId], data)
		server.GlobalLogCountState[peerId] ++
	} else {
		server.GlobalLogState[peerId] = make([]synceddoc.OperationRequest, 1)
		server.GlobalLogState[peerId][0] = data
		server.GlobalLogCountState[peerId] = 1
	}
}

func (server *Server) GetLog() (map[utils.UUID][]synceddoc.OperationRequest, map[utils.UUID]int) {
	return server.GlobalLogState, server.GlobalLogCountState
}