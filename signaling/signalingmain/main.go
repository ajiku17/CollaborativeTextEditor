package main

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/signaling"
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

	s := signaling.NewServer()

	err := http.ListenAndServe(port, s)
	if err != nil {
		fmt.Println("error: ", err)
	}
}
