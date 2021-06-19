package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ajiku17/CollaborativeTextEditor/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type RemoteDocument struct {
	url    string
	httpClient *http.Client
}

type Request struct {
	Id string `json:"id"`
	Site string `json:"site"`
	Position string `json:"position"`
	Value string `json:"value"`
}

type Response struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
}

//site???
func (doc *RemoteDocument)InsertAtPosition(position crdt.Position, val string) {	
	site := val[strings.Index(val, ":") + 1:]
	jsonStr := utils.ToJson(Request{"1", site, crdt.BasicPositionToString(position.(crdt.BasicPosition)), val})
	doc.sendRequest(bytes.NewBuffer(jsonStr), "Insert")
}

func (doc *RemoteDocument)DeleteAtPosition(position crdt.Position) {
	site := "1"
	jsonStr := utils.ToJson(Request{"1", site, crdt.BasicPositionToString(position.(crdt.BasicPosition)), ""})
	doc.sendRequest(bytes.NewBuffer(jsonStr), "Delete")
}

func (doc *RemoteDocument) sendRequest(data *bytes.Buffer, methodName string) int {

	req, err := http.NewRequest("POST", fmt.Sprintf(doc.url + methodName), data)
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := doc.httpClient.Do(req)
	if err != nil {
		return 0
	}

	handleResponse(resp)
	return 0
}

func handleResponse(resp *http.Response) {
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	responseObject := utils.FromJson(bodyBytes, Response{})
	fmt.Printf("Response - %+v\n", responseObject)
}