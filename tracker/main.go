package main

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/tracker/server"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Usage: server [port]")
		return
	}

	port := ":" + args[1]

	err := server.Start(port)
	if err != nil {
		fmt.Println("error: ", err)
	}
}
