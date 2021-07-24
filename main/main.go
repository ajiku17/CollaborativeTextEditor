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
	site2 := utils.GenerateNewUUID()
	fmt.Println("Site2 - ", site2)
	doc2 := synceddoc.New(string(site2))
	manager2 := network.NewDocumentManager(site2, &doc2)
	doc2.ConnectSignals(onChangeListener, nil, nil)
	manager2.Start()
	time.Sleep(2*time.Second)


	doc1.LocalInsert(doc1.GetDocument().Length(), "H")
	time.Sleep(2*time.Second)
	doc2.LocalInsert(doc2.GetDocument().Length(), "E")
	time.Sleep(2*time.Second)
	doc1.LocalInsert(doc1.GetDocument().Length(), "L")
	time.Sleep(2*time.Second)
	fmt.Printf("Document for site %s is : %s\n", doc1.GetSiteID(), doc1.GetDocument().ToString())
	fmt.Printf("Document for site %s is : %s\n", doc2.GetSiteID(), doc2.GetDocument().ToString())

	doc2.LocalDelete(0)
	time.Sleep(2*time.Second)

	manager1.Kill()
	manager2.Kill()

	fmt.Printf("Document for site %s is : %s\n", doc1.GetSiteID(), doc1.GetDocument().ToString())
	fmt.Printf("Document for site %s is : %s\n", doc2.GetSiteID(), doc2.GetDocument().ToString())
	fmt.Printf("Log for site %s is : %s\n", doc1.GetSiteID(), doc1.GetLogs())
	fmt.Printf("Log for site %s is : %s\n", doc2.GetSiteID(), doc2.GetLogs())
	log, counts := server.GetLog()
	fmt.Printf("Global Log %s", log)
	fmt.Printf("Global Log %s", counts)
}

// Example: d.onChange(MESSAGE_INSERT, MessageInsert{Index: index, Value: val}, aux)
func onChangeListener(changeName string, change interface {}, aux interface{}) {
}

func peerConnectedListener(peerId utils.UUID, cursorPosition int, aux interface{}) {}
func peerDisconnectedListener(peerId utils.UUID, aux interface{}) {}