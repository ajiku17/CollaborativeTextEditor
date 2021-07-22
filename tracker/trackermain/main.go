package main

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/tracker"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Usage: server [port]")
		return
	}

	port := ":" + args[1]

	tr := tracker.NewHttpTracker()

	err := http.ListenAndServe(port, tr)
	if err != nil {
		fmt.Println("error: ", err)
	}
}
