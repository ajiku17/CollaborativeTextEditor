package main

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/core/server"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"time"
)

func main() {
	server := server.NewServer()
	go server.Listen()

	site1 := utils.GenerateNewUUID()
	fmt.Println("Site1 - ", site1)
	doc1 := synceddoc.New(string(site1))
	manager1 := network.NewDocumentManager(site1, &doc1)
	doc1.ConnectSignals(onChangeListener, nil, nil)
	manager1.Start()


	time.Sleep(100)
	doc1.LocalInsert(doc1.GetDocument().Length(), "H")
	doc1.LocalInsert(doc1.GetDocument().Length(), "E")


	site2 := utils.GenerateNewUUID()
	fmt.Println("Site2 - ", site2)
	doc2 := synceddoc.New(string(site2))
	manager2 := network.NewDocumentManager(site2, &doc2)
	doc2.ConnectSignals(onChangeListener, nil, nil)
	manager2.Start()
	//doc1.LocalInsert(doc1.GetDocument().Length(), "L")
	time.Sleep(10000000)

	fmt.Printf("Document for site %s is : %s\n", doc1.GetSiteID(), doc1.GetDocument().ToString())
	fmt.Printf("Document for site %s is : %s\n", doc2.GetSiteID(), doc2.GetDocument().ToString())
}

// Example: d.onChange(MESSAGE_INSERT, MessageInsert{Index: index, Value: val}, aux)
func onChangeListener(changeName string, change interface {}, aux interface{}) {
	switch changeName {
	case synceddoc.MESSAGE_INSERT:
		fmt.Println("Insert")
	case synceddoc.MESSAGE_DELETE:
		fmt.Println("Delete")
	case synceddoc.MESSAGE_PEER_CURSOR:
		fmt.Println("Cursor")
	}
}

func peerConnectedListener(peerId utils.UUID, cursorPosition int, aux interface{}) {}
func peerDisconnectedListener(peerId utils.UUID, aux interface{}) {}