package crdt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type DocumentUpdateManager struct {
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

func (manager *DocumentUpdateManager) Insert(position Position, val string, site int) {
	jsonStr := utils.ToJson(Request{"1", strconv.Itoa(site), BasicPositionToString(position.(BasicPosition)), val})
	manager.sendRequest(bytes.NewBuffer(jsonStr), "Insert")
}


func (manager *DocumentUpdateManager) Delete(position Position, site int) {
	jsonStr := utils.ToJson(Request{"1", strconv.Itoa(site), BasicPositionToString(position.(BasicPosition)), ""})
	manager.sendRequest(bytes.NewBuffer(jsonStr), "Delete")
}

func (manager *DocumentUpdateManager)sendRequest(data *bytes.Buffer, methodName string) int {
	socket, err := net.Dial("tcp", manager.url)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	
	for {
		text := "\n" + methodName
		socket.Write(append(data.Bytes(), []byte(text)...))
		fmt.Printf("text - %s\n", text)

		var received []byte
		socket.Read(received)
		fmt.Printf("->: %b\n", received)
		if strings.TrimSpace(string(text)) == "" {
				fmt.Println("TCP client exiting...")
				return 0
		}
	}
}

func  (manager *DocumentUpdateManager)sendHttpRequest(data *bytes.Buffer, methodName string) int {
	req, err := http.NewRequest("POST", fmt.Sprintf(manager.url+methodName), data)
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := manager.httpClient.Do(req)
	if err != nil {
		return 0
	}

	handleHttpResponse(resp)
	return 0
}

func handleHttpResponse(resp *http.Response) {
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	responseObject := utils.FromJson(bodyBytes, Response{})
	fmt.Printf("Response - %+v\n", responseObject)
}