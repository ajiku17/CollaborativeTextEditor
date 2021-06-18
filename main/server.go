package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ajiku17/CollaborativeTextEditor/client"
	"github.com/ajiku17/CollaborativeTextEditor/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Server struct {
	Clients []*(client.Client)
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
	var request map[string]interface{} = utils.FromJson(jsonBytes, client.Request{}).(map[string]interface{})

	site :=  int(request["site"].(float64))
	position := crdt.ToBasicPosition(request["position"].(string))
	value := request["value"].(string)

	for _, client := range server.Clients {
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

	var request map[string]interface{} = utils.FromJson(jsonBytes, client.Request{}).(map[string]interface{})

	site :=  int(request["site"].(float64))
	position := crdt.ToBasicPosition(request["position"].(string))

	for _, client := range server.Clients {
		if client.GetSite() != site {
			client.DeleteAtPosition(position)
		}
	}
}

func (server *Server)HandleRequests() {
    http.HandleFunc("/Insert", insert)
    http.HandleFunc("/Delete", delete)
    log.Fatal(http.ListenAndServe(":8081", nil))
}

func (server *Server)ConnectWithClient(client *client.Client) {
	server.Clients = append(server.Clients, client)
}