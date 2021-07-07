package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Tracker struct {
	table map[string] []string
	mu    sync.Mutex
}

var tracker = new(Tracker)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	arguments := r.URL.Query()

	docId, ok := arguments["doc"]
	if !ok || len(docId) < 1 {
		http.Error(w, errors.New("must provide document id 'doc'").Error(), http.StatusPreconditionFailed)
		return
	}

	addr, ok := arguments["addr"]
	if !ok || len(addr) < 1 {
		http.Error(w, errors.New("must provide address 'addr'").Error(), http.StatusPreconditionFailed)
		return
	}

	remoteAddr := addr[0]

	tracker.mu.Lock()
	peers, ok := tracker.table[docId[0]]
	if !ok {
		peers = []string{}
	}

	tracker.table[docId[0]] = append(peers, remoteAddr)
	tracker.mu.Unlock()
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	arguments := r.URL.Query()

	docId, ok := arguments["doc"]
	if !ok || len(docId) < 1 {
		http.Error(w, errors.New("must provide document id 'doc'").Error(), http.StatusPreconditionFailed)
		return
	}

	tracker.mu.Lock()
	
	peers, ok := tracker.table[docId[0]]
	if !ok {
		peers = []string{}
	}

	tracker.mu.Unlock()

	peersJson, _ := json.Marshal(peers)
	w.Header().Set("Content-Type", "application/json")

	_, err := w.Write(peersJson)
	if err != nil {
		fmt.Println("an error occurred while writing", err)
	}
}

func initState() {
	tracker.table = make(map[string] []string)
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Usage: server [port]")
		return
	}

	port := ":" + args[1]

	initState()

	http.HandleFunc("/register/", registerHandler)
	http.HandleFunc("/get/", getHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
