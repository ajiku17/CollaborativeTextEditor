package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/crdt"
	"github.com/utils"
)

type Server struct {
	Clients []*crdt.Client
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

	var request map[string]interface{}
	request = utils.FromJson(jsonBytes, crdt.Request{}).(map[string]interface{})

	site :=  int(request["site"].(float64))
	position := crdt.ToPosition(request["position"].(string))
	value := request["value"].(string)

	// fmt.Printf("server client length is %d\n", len(server.Clients))
	// fmt.Printf("curr site is %d", site)
	for _, client := range server.Clients {
		if client.GetSite() != site {
			// fmt.Printf("send to positiion %d value %s\n", client.GetSite(), value)
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

	var request map[string]interface{}
	request = utils.FromJson(jsonBytes, crdt.Request{}).(map[string]interface{})

	site :=  int(request["site"].(float64))
	position := crdt.ToPosition(request["position"].(string))

	// fmt.Printf("server client length is %d\n", len(server.Clients))
	// fmt.Printf("curr site is %d", site)
	for _, client := range server.Clients {
		if client.GetSite() != site {
			// fmt.Printf("send to positiion %d value %s\n", client.GetSite(), value)
			client.DeleteAtPosition(position)
		}
	}
}

func (server *Server)HandleRequests() {
    http.HandleFunc("/Insert", insert)
    http.HandleFunc("/Delete", delete)
    log.Fatal(http.ListenAndServe(":8081", nil))
}

func (server *Server)ConnectWithClient(client *crdt.Client) {
	server.Clients = append(server.Clients, client)
}