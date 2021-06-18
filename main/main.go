package main

import (
	"log"

	"github.com/ajiku17/CollaborativeTextEditor/client"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server := NewServer()
	go server.HandleRequests()

	client1 := client.NewClient(1)
	client2 := client.NewClient(2)
	server.ConnectWithClient(client1)
	server.ConnectWithClient(client2)

	client1.Insert("H", 0)
	client2.Insert("e", 1)
	
	client1.Insert("l", 2)
	client1.Insert("l", 3)
	client2.Insert("o", 4)
	client2.Insert(" ", 5)
	client2.Insert("W", 6)
	client1.Insert("o", 7)
	client1.Insert("r", 8)
	client1.Insert("l", 9)
	client1.Delete(9)
	client2.Insert("d", 9)

	client1.PrintDocument()
	client2.PrintDocument()
}
