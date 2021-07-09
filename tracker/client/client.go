package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	serverURL string
}

func New(url string) *Client {
	client := new(Client)
	client.serverURL = url

	return client
}

// Register registers specified address as an editor of the specified document
func (c *Client) Register(docuemntId string, address string) {
	_, err := http.Get(c.serverURL + "/register?doc=" + docuemntId + "&addr=" + address)
	if err != nil {
		fmt.Printf("error occured: %v", err)
		return
	}
}

func parseGet(r io.Reader) []string {
	peerList := make([]string, 0)

	rspData, err := ioutil.ReadAll(r)
	if err != nil {
		return peerList
	}

	if len(rspData) == 0 {
		fmt.Println("error while retrieving data", err)
		return peerList
	}

	err = json.Unmarshal(rspData, &peerList)
	if err != nil && err != io.EOF {
		fmt.Println("error while parsing", err)
	}

	return peerList
}

// Get retrieves a list currently connected peers
func (c *Client) Get(docuemntId string) []string {
	rsp, err := http.Get(c.serverURL + "/get?doc=" + docuemntId)
	if err != nil {
		fmt.Printf("error occured: %v", err)
		return []string{}
	}

	return parseGet(rsp.Body)
}
