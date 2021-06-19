package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

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
	return server
}

func insert(w http.ResponseWriter, r *http.Request){
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	print("%s\n", string(jsonBytes))

	
	// Send request to all clients
	var request map[string]interface{} = utils.FromJson(jsonBytes, crdt.Request{}).(map[string]interface{})

	site, err :=  strconv.Atoi(request["site"].(string))
	position := crdt.ToBasicPosition(request["position"].(string))
	value := request["value"].(string)

	for _, client := range server.SyncedDocuments {
		if client.GetSite() != site {
			client.InsertAtPosition(position, value)
		}
	}
}


func delete(w http.ResponseWriter, r *http.Request){
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	print("%s\n", string(jsonBytes))

	
	// Send request to all clients

	var request map[string]interface{} = utils.FromJson(jsonBytes, crdt.Request{}).(map[string]interface{})

	site :=  request["site"]
	position := crdt.ToBasicPosition(request["position"].(string))

	for _, client := range server.SyncedDocuments {
		if client.GetSite() != site {
			client.DeleteAtPosition(position)
		}
	}
}

// func (server *Server)HandleRequests() {
//     http.HandleFunc("/Insert", insert)
//     http.HandleFunc("/Delete", delete)
//     log.Fatal(http.ListenAndServe(":8081", nil))
// }

func (server *Server)ConnectWithClient(doc *crdt.SyncedDocument) {
	server.SyncedDocuments = append(server.SyncedDocuments, doc)
}

func (server *Server) HandleRequests(mu *sync.Mutex) {
	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
			fmt.Println(err)
			mu.Unlock()
			return
	}
	mu.Unlock()

	defer listener.Close()

	socket, err := listener.Accept()
	if err != nil {
			fmt.Println(err)
			return
	}

	// fmt.Printf("%s\n", socket)
	for {
			var receivedMessage []byte
			socket.Read(receivedMessage)

			if err != nil {
					fmt.Println(err)
					return
			}
			fmt.Printf("RECEIVED - %b\n", receivedMessage)
			if strings.TrimSpace(string(receivedMessage)) == "" {
					fmt.Println("Exiting TCP server!")
					return
			}
			
			sendMessage := "Hi"
			socket.Write([]byte(sendMessage))
	}
}