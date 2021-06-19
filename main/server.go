package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ajiku17/CollaborativeTextEditor/client"
	"github.com/ajiku17/CollaborativeTextEditor/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Server struct {
	SynchedDocuments []*(client.SynchedDocument)
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

	site, err :=  strconv.Atoi(request["site"].(string))
	position := crdt.ToBasicPosition(request["position"].(string))
	value := request["value"].(string)

	for _, client := range server.SynchedDocuments {
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

	site :=  request["site"]
	position := crdt.ToBasicPosition(request["position"].(string))

	for _, client := range server.SynchedDocuments {
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

func (server *Server)ConnectWithClient(doc *client.SynchedDocument) {
	server.SynchedDocuments = append(server.SynchedDocuments, doc)
}