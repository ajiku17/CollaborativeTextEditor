package main

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"syscall/js"
)

var docManager *DocumentManager

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

func buildPeerConnectedCallback(peerConnectedCallback js.Value) synceddoc.PeerConnectedListener {
	return func (peerId utils.UUID, cursorPosition int, aux interface{}) {
		peerConnectedCallback.Invoke(string(peerId), cursorPosition)
	}
}

func buildPeerDisconnectedCallback(peerDisconnectedCallback js.Value) synceddoc.PeerDisconnectedListener {
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

	doc, err := synceddoc.Open(documentId)
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	doc.ConnectSignals(buildChangeCallback(changeCallback),
		buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))
	docId := docManager.PutDocument(doc)
	initCallback.Invoke(doc.ToString())

	return docId
}

func DocumentDeserialize(this js.Value, i []js.Value) interface {} {
	byteArray := i[0]
	initCallback := i[1]
	changeCallback := i[2]
	peerConnectedCallback := i[3]
	peerDisconnectedCallback := i[4]

	serialized := make([]byte, byteArray.Get("length").Int())

	js.CopyBytesToGo(serialized, i[0])

	doc, err := synceddoc.Load(serialized)
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	doc.ConnectSignals(buildChangeCallback(changeCallback),
		buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))
	docId := docManager.PutDocument(doc)

	// call init callback
	initCallback.Invoke(doc.ToString())

	return string(docId)
}

func DocumentNew(this js.Value, i []js.Value) interface {} {
	changeCallback := i[0]
	peerConnectedCallback := i[1]
	peerDisconnectedCallback := i[2]

	doc := synceddoc.New()

	doc.ConnectSignals(buildChangeCallback(changeCallback),
		buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))

	docId := docManager.PutDocument(doc)

	return string(docId)
}

func DocumentClose(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	doc.Close()

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

	doc.LocalInsert(index, value)

	fmt.Printf("Calling DocumentInsertAt on docId: %v index: %d, value: %s after insert: %s\n", docId, index, value, doc.ToString())
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

	doc.LocalDelete(index)

	fmt.Printf("Calling DocumentDeleteAt on docId: %v index: %d, after delete: %s\n", docId, index, doc.ToString())
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

	doc.SetCursor(index)
	fmt.Printf("Calling DocumentChangeCursor on docId: %v index: %d doc: %s\n", docId, index, doc.ToString())
	return nil
}

func DocumentSerialize(this js.Value, i []js.Value) interface {} {
	docId := i[0].String()

	doc, err := docManager.GetDocument(DocumentID(docId))

	if err != nil {
		return []interface{}{-1}
	}

	serialized, err := doc.Serialize()
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

func SetPeerConnectedListener(doc synceddoc.Document, callback js.Value) {
	doc.SetPeerConnectedListener(buildPeerConnectedCallback(callback))
}

func SetPeerDisconnectedListener(doc synceddoc.Document, callback js.Value) {
	doc.SetPeerDisconnectedListener(buildPeerDisconnectedCallback(callback))
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
	case "onChange": SetChangeListener(doc, callback)
	case "onPeerConnect": SetPeerConnectedListener(doc, callback)
	case "onPeerDisconnect": SetPeerDisconnectedListener(doc, callback)
	}

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
}

func main() {
	c := make(chan struct{}, 0)

	docManager = NewDocumentManager()
	registerCallbacks()
	fmt.Println("Callbacks registered")

	<- c // prevent main from terminating
}