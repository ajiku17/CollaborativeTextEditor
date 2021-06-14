package main

import (
	"github.com/crdt"
)



func main() {
	server := NewServer()
	go server.HandleRequests()

	// doc.InsertAt("H", 0, 1)
	client1 := crdt.NewClient(1)
	client2 := crdt.NewClient(2)
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
