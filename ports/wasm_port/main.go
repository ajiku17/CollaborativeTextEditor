package main

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/tracker"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"syscall/js"
)

const TRACKER_URL = "http://127.0.0.1:9090"
const SIGNALING_URL = "http://127.0.0.1:9999"

var docManager    *DocumentManager
var siteId         utils.UUID

func buildChangeCallback(changeCallback js.Value) synceddoc.ChangeListener {
	return func (changeName string, change interface {}, aux interface{}) {
		changeObj := make(map[string]interface{})

		changeObj["changeName"] = changeName

		switch change.(type) {
		case synceddoc.MessageInsert:
			changeObj["index"] = change.(synceddoc.MessageInsert).Index
			changeObj["value"] = change.(synceddoc.MessageInsert).Value
		case synceddoc.MessageDelete:
			changeObj["index"] = change.(synceddoc.MessageDelete).Index
		case synceddoc.MessagePeerCursor:
			changeObj["peerId"] = change.(synceddoc.MessagePeerCursor).PeerID
			changeObj["cursorPos"] = change.(synceddoc.MessagePeerCursor).CursorPosition
		}

		changeCallback.Invoke(changeName, changeObj)
	}
}

func buildPeerConnectedCallback(peerConnectedCallback js.Value) network.PeerConnectedListener {
	return func (peerId utils.UUID, cursorPosition int, aux interface{}) {
		peerConnectedCallback.Invoke(string(peerId), cursorPosition)
	}
}

func buildPeerDisconnectedCallback(peerDisconnectedCallback js.Value) network.PeerDisconnectedListener {
	return func (peerId utils.UUID, aux interface{}) {
		peerDisconnectedCallback.Invoke(string(peerId))
	}
}

func DocumentOpen(this js.Value, i []js.Value) interface {} {
	documentId := i[0].String()
	initCallback := i[1]
	changeCallback := i[2]
	peerConnectedCallback := i[3]
	peerDisconnectedCallback := i[4]
	resolve := i[5]

	siteId = utils.GenerateNewUUID()
	doc, err := synceddoc.Open(string(siteId), documentId)
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	trackerC := tracker.NewClient(TRACKER_URL)
	manager := network.NewP2PManager(siteId, doc, SIGNALING_URL, trackerC)

	doc.ConnectSignals(buildChangeCallback(changeCallback))
	manager.ConnectSignals(buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))
	docManager.PutDocument(Document{
		Doc:        doc,
		NetManager: manager,
	})

	// call js init callback
	initCallback.Invoke(doc.ToString(), string(siteId))

	go func () {
		manager.Start()

		resolve.Invoke(string(doc.GetID()))
	}()

	return nil
}

func DocumentDeserialize(this js.Value, i []js.Value) interface {} {
	byteArray := i[0]
	initCallback := i[1]
	changeCallback := i[2]
	peerConnectedCallback := i[3]
	peerDisconnectedCallback := i[4]
	resolve := i[5]

	serialized := make([]byte, byteArray.Get("length").Int())

	js.CopyBytesToGo(serialized, i[0])

	siteId = utils.GenerateNewUUID()
	doc, err := synceddoc.Load(string(siteId), serialized)
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	trackerC := tracker.NewClient(TRACKER_URL)
	manager := network.NewP2PManager(siteId, doc, SIGNALING_URL, trackerC)

	doc.ConnectSignals(buildChangeCallback(changeCallback))
	manager.ConnectSignals(buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))
	docManager.PutDocument(Document{
		Doc:        doc,
		NetManager: manager,
	})

	// call init callback
	initCallback.Invoke(doc.ToString(), string(siteId))

	go func () {
		manager.Start()

		resolve.Invoke(string(doc.GetID()))
	}()

	return nil
}

func DocumentNew(this js.Value, i []js.Value) interface {} {
	changeCallback := i[0]
	peerConnectedCallback := i[1]
	peerDisconnectedCallback := i[2]
	initCallback := i[3]
	resolve := i[4]

	siteId = utils.GenerateNewUUID()
	doc := synceddoc.New(string(siteId))
	trackerC := tracker.NewClient(TRACKER_URL)
	manager := network.NewP2PManager(siteId, doc, SIGNALING_URL, trackerC)

	doc.ConnectSignals(buildChangeCallback(changeCallback))
	manager.ConnectSignals(buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))

	docManager.PutDocument(Document{
		Doc:        doc,
		NetManager: manager,
	})

	initCallback.Invoke("", string(siteId))

	go func () {
		manager.Start()

		resolve.Invoke(string(doc.GetID()))
	}()

	return nil
}

func DocumentClose(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	fmt.Println("document CLOSING")
	go func() {
		doc.Doc.Close()
		doc.NetManager.Stop()
	} ()

	fmt.Println("document CLOSED")
	docManager.RemoveDocument(DocumentID(docId))

	return nil
}

func DocumentInsertAt(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()
	value := i[1].String()
	index := i[2].Int()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		fmt.Println("error: ", err)
		return -1
	}

	doc.Doc.LocalInsert(index, value)

	fmt.Printf("Calling DocumentInsertAt on docId: %v index: %d, value: %s after insert: %s\n", docId, index, value, doc.Doc.ToString())
	return nil
}

func DocumentDeleteAt(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()
	index := i[1].Int()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		fmt.Println(err)
		return -1
	}

	doc.Doc.LocalDelete(index)

	fmt.Printf("Calling DocumentDeleteAt on docId: %v index: %d, after delete: %s\n", docId, index, doc.Doc.ToString())
	return nil
}

func DocumentChangeCursor(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()
	index := i[1].Int()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		fmt.Println(err)
		return -1
	}

	doc.Doc.SetCursor(index)
	fmt.Printf("Calling DocumentChangeCursor on docId: %v index: %d doc: %s\n", docId, index, doc.Doc.ToString())
	return nil
}

func DocumentSerialize(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		return []interface{}{-1}
	}

	serialized, err := doc.Doc.Serialize()
	if err != nil {
		return -1
	}

	res := make([]interface{}, len(serialized))

	for i, val := range serialized {
		res[i] = val
	}

	return res
}

func SetChangeListener(doc synceddoc.Document, callback js.Value) {
	doc.SetChangeListener(buildChangeCallback(callback))
}

func SetPeerConnectedListener(manager network.Manager, callback js.Value) {
	manager.OnPeerConnect(buildPeerConnectedCallback(callback))
}

func SetPeerDisconnectedListener(manager network.Manager, callback js.Value) {
	manager.OnPeerDisconnect(buildPeerDisconnectedCallback(callback))
}

func DocumentSetListener(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()
	callbackName := i[1].String()
	callback := i[2]

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		return -1
	}

	switch callbackName {
	case "onChange": SetChangeListener(doc.Doc, callback)
	case "onPeerConnect": SetPeerConnectedListener(doc.NetManager, callback)
	case "onPeerDisconnect": SetPeerDisconnectedListener(doc.NetManager, callback)
	}

	return nil
}

func DocumentDisconnect(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	fmt.Println("Document Disconnect")

	go func () {
		doc.NetManager.Stop()
	}()

	return nil
}

func DocumentReconnect(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	go func () {
		fmt.Println("Reconnecting ", doc.Doc)
		doc.NetManager.Start()
	}()

	return nil
}

func registerCallbacks() {
	js.Global().Set("DocumentOpen", js.FuncOf(DocumentOpen))
	js.Global().Set("DocumentDeserialize", js.FuncOf(DocumentDeserialize))
	js.Global().Set("DocumentNew", js.FuncOf(DocumentNew))
	js.Global().Set("DocumentClose", js.FuncOf(DocumentClose))
	js.Global().Set("DocumentInsertAt", js.FuncOf(DocumentInsertAt))
	js.Global().Set("DocumentDeleteAt", js.FuncOf(DocumentDeleteAt))
	js.Global().Set("DocumentChangeCursor", js.FuncOf(DocumentChangeCursor))
	js.Global().Set("DocumentSerialize", js.FuncOf(DocumentSerialize))
	js.Global().Set("DocumentSetListener", js.FuncOf(DocumentSetListener))
	js.Global().Set("DocumentDisconnect", js.FuncOf(DocumentDisconnect))
	js.Global().Set("DocumentReconnect", js.FuncOf(DocumentReconnect))
}

func main() {
	c := make(chan struct{}, 0)

	docManager = NewDocumentManager()
	registerCallbacks()
	fmt.Println("Callbacks registered")

	<- c // prevent main from terminating
}