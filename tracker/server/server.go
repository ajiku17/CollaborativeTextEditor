package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/tracker/table"
	"net/http"
)

type HttpTracker struct {
	Table *table.Table
}

func (t *HttpTracker) registerHandler(w http.ResponseWriter, r *http.Request) {
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

	t.Table.Register(docId[0], addr[0])
}

func (t *HttpTracker) getHandler(w http.ResponseWriter, r *http.Request) {
	arguments := r.URL.Query()

	docId, ok := arguments["doc"]
	if !ok || len(docId) < 1 {
		http.Error(w, errors.New("must provide document id 'doc'").Error(), http.StatusPreconditionFailed)
		return
	}

	peersJson, err := json.Marshal(t.Table.Get(docId[0]))
	if err != nil {
		fmt.Println("an error occurred while marshaling", err)
	}
	w.Header().Set("Content-Type", "application/json")

	peerList := make([]string, 0)
	err = json.Unmarshal(peersJson, &peerList)
	if err != nil {
		fmt.Println("an error occurred while unmarshaling", err)
	}

	_, err = w.Write(peersJson)
	if err != nil {
		fmt.Println("an error occurred while writing", err)
	}
}

func (t *HttpTracker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path  {
	case "/get":
		t.getHandler(w, r)
	case "/register":
		t.registerHandler(w, r)
	}
}

func NewHttpTracker() *HttpTracker {
	t := new(HttpTracker)

	t.Table = table.New()
	return t
}

func Start (port string) error {
	tracker := NewHttpTracker()

	http.Handle("/", tracker)

	return http.ListenAndServe(port, nil)
}
