package main

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
	"syscall/js"
)

var docManager *DocumentManager

func buildChangeCallback(changeCallback js.Value) synceddoc.ChangeListener {
	return func (changeName string, change interface {}) {
		changeObj := make(map[string]interface{})

		changeObj["changeName"] = changeName

		switch change.(type) {
		case synceddoc.ChangeInsert:
			changeObj["index"] = change.(synceddoc.ChangeInsert).Index
			changeObj["value"] = change.(synceddoc.ChangeInsert).Value
		case synceddoc.ChangeDelete:
			changeObj["index"] = change.(synceddoc.ChangeDelete).Index
		case synceddoc.ChangePeerCursor:
			changeObj["peerId"] = change.(synceddoc.ChangePeerCursor).PeerID
			changeObj["cursorPos"] = change.(synceddoc.ChangePeerCursor).CursorPosition
		}

		changeCallback.Invoke(changeName, changeObj)
	}
}

func buildPeerConnectedCallback(peerConnectedCallback js.Value) synceddoc.PeerConnectedListener {
	return func (peerId utils.UUID, cursorPosition int) {
		peerConnectedCallback.Invoke(string(peerId), cursorPosition)
	}
}

func buildPeerDisconnectedCallback(peerDisconnectedCallback js.Value) synceddoc.PeerDisconnectedListener {
	return func (peerId utils.UUID) {
		peerDisconnectedCallback.Invoke(string(peerId))
	}
}

func DocumentOpen(this js.Value, i []js.Value) interface {} {
	documentId := i[0].String()
	initCallback := i[1]
	changeCallback := i[2]
	peerConnectedCallback := i[3]
	peerDisconnectedCallback := i[4]

	doc, err := synceddoc.Open(documentId,
		buildChangeCallback(changeCallback),
		buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	fd := docManager.PutDocument(doc)
	initCallback.Invoke(doc.ToString())

	return fd
}

func DocumentDeserialize(this js.Value, i []js.Value) interface {} {
	byteArray := i[0]
	initCallback := i[1]
	changeCallback := i[2]
	peerConnectedCallback := i[3]
	peerDisconnectedCallback := i[4]

	serialized := make([]byte, byteArray.Get("length").Int())

	js.CopyBytesToGo(serialized, i[0])

	doc, err := synceddoc.Load(serialized,
		buildChangeCallback(changeCallback),
		buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	fd := docManager.PutDocument(doc)

	// call init callback
	initCallback.Invoke(doc.ToString())

	return int(fd)
}

func DocumentNew(this js.Value, i []js.Value) interface {} {
	changeCallback := i[0]
	peerConnectedCallback := i[1]
	peerDisconnectedCallback := i[2]

	doc := synceddoc.New(buildChangeCallback(changeCallback),
		buildPeerConnectedCallback(peerConnectedCallback),
		buildPeerDisconnectedCallback(peerDisconnectedCallback))

	fd := docManager.PutDocument(doc)

	return int(fd)
}

func DocumentClose(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()

	doc, err := docManager.GetDocument(FileDescriptor(fd))

	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}

	doc.Close()

	docManager.RemoveDocument(FileDescriptor(fd))

	return nil
}

func DocumentInsertAt(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()
	value := i[1].String()
	index := i[2].Int()

	doc, err := docManager.GetDocument(FileDescriptor(fd))

	if err != nil {
		fmt.Println("error: ", err)
		return -1
	}

	doc.InsertAtIndex(index, value)

	fmt.Printf("Calling DocumentInsertAt on fd: %v index: %d, value: %s after insert: %s\n", fd, index, value, doc.ToString())
	return nil
}

func DocumentDeleteAt(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()
	index := i[1].Int()

	doc, err := docManager.GetDocument(FileDescriptor(fd))

	if err != nil {
		fmt.Println(err)
		return -1
	}

	doc.DeleteAtIndex(index)

	fmt.Printf("Calling DocumentDeleteAt on fd: %v index: %d, after delete: %s\n", fd, index, doc.ToString())
	return nil
}

func DocumentChangeCursor(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()
	index := i[1].Int()

	doc, err := docManager.GetDocument(FileDescriptor(fd))

	if err != nil {
		fmt.Println(err)
		return -1
	}

	doc.SetCursor(index)
	fmt.Printf("Calling DocumentChangeCursor on fd: %v index: %d doc: %s\n", fd, index, doc.ToString())
	return nil
}

func DocumentSerialize(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()

	doc, err := docManager.GetDocument(FileDescriptor(fd))

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
	fd := i[0].Int()
	callbackName := i[1].String()
	callback := i[2]

	doc, err := docManager.GetDocument(FileDescriptor(fd))

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