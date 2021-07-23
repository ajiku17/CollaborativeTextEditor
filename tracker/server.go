package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HttpTracker struct {
	Table *Table

	serveMux http.ServeMux
}

func (t *HttpTracker) registerHandler(w http.ResponseWriter, r *http.Request) {
	arguments := r.URL.Query()

	docId, ok := arguments["doc"]
	if !ok || len(docId) < 1 {
		http.Error(w, errors.New("must provide document id 'doc'").Error(), http.StatusPreconditionFailed)
		return
	}

	addr, ok := arguments["peerid"]
	if !ok || len(addr) < 1 {
		http.Error(w, errors.New("must provide peer id 'peerid'").Error(), http.StatusPreconditionFailed)
		return
	}

	t.Table.Register(docId[0], addr[0])

	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = w.Write(peersJson)
	if err != nil {
		fmt.Println("an error occurred while writing", err)
	}
}

func (t *HttpTracker) registerGetHandler(w http.ResponseWriter, r *http.Request) {
	arguments := r.URL.Query()

	docId, ok := arguments["doc"]
	if !ok || len(docId) < 1 {
		http.Error(w, errors.New("must provide document id 'doc'").Error(), http.StatusPreconditionFailed)
		return
	}

	addr, ok := arguments["peerid"]
	if !ok || len(addr) < 1 {
		http.Error(w, errors.New("must provide peer id 'peerid'").Error(), http.StatusPreconditionFailed)
		return
	}

	peersJson, err := json.Marshal(t.Table.RegisterAndGet(docId[0], addr[0]))
	if err != nil {
		fmt.Println("an error occurred while marshaling", err)
	}
	w.Header().Set("Content-Type", "application/json")

	peerList := make([]string, 0)
	err = json.Unmarshal(peersJson, &peerList)
	if err != nil {
		fmt.Println("an error occurred while unmarshaling", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = w.Write(peersJson)
	if err != nil {
		fmt.Println("an error occurred while writing", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func NewHttpTracker() *HttpTracker {
	t := new(HttpTracker)

	t.Table = NewTable()

	t.serveMux.HandleFunc("/get", t.getHandler)
	t.serveMux.HandleFunc("/register", t.registerHandler)
	t.serveMux.HandleFunc("/register-get", t.registerGetHandler)

	return t
}

func (s *HttpTracker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func Start (port string) error {
	tracker := NewHttpTracker()

	http.Handle("/", tracker)

	return http.ListenAndServe(port, nil)
}
