package crdt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/utils"
)

type ClientServer struct {
	url    string
	httpClient *http.Client
}

type Request struct {
	Id string `json:"id"`
	Site int `json:"site"`
	Position string `json:"position"`
	Value string `json:"value"`
}

type Response struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
}

func (client *ClientServer) SendInsertRequest(position Position, val string, site int) int {
	jsonStr := utils.ToJson(Request{"1", site, position.ToString(), val})
	return client.sendRequest(bytes.NewBuffer(jsonStr), "Insert")
}

func (client *ClientServer) SendDeleteRequest(position Position, site int) int {
	jsonStr := utils.ToJson(Request{"1", site, position.ToString(), ""})
	return client.sendRequest(bytes.NewBuffer(jsonStr), "Delete")
}

func (client *ClientServer) sendRequest(data *bytes.Buffer, methodName string) int {

	req, err := http.NewRequest("POST", fmt.Sprintf(client.url + methodName), data)
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := client.httpClient.Do(req)
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