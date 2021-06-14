package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/crdt"
)

type Client struct {
	site int
	clientServer *ClientServer
	document     *crdt.Document
}

func NewClient(site int) *Client {
	server_url := "http://localhost:8081/"
	client := Client{site, &ClientServer{server_url, &http.Client{Timeout: 5 * time.Minute}}, crdt.NewDocument()}
	return &client
}

func (client *Client) GetSite() int {
	return client.site
}

func (client *Client) Insert(val string, index int) {
	position := client.document.InsertAt(val, index, client.site)
	client.clientServer.SendInsertRequest(position, val, client.site)
}

func (client *Client) InsertAtPosition(pos crdt.Position, val string) {
	client.document.InsertAtPos(pos, val)
	// TODO: send server an acknowledgement request
}

func (client *Client) Delete(index int) {
	position := client.document.DeleteAt(index)
	client.clientServer.SendDeleteRequest(position, client.site)
}

func (client *Client) DeleteAtPosition(pos crdt.Position) {
	client.document.DeleteAtPos(pos)
	// TODO: send server an acknowledgement request
}


func (client *Client)PrintDocument() {
	fmt.Printf("Document for client site N %d is : \n", client.site)
	fmt.Println(client.document.ToString())
}