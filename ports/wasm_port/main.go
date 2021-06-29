package main

import (
	"fmt"
	"github.com/ajiku17/CollaborativeTextEditor/crdt"
	"syscall/js"
)

var docManager *DocumentManager

func DocumentOpen(this js.Value, i []js.Value) interface {} {
	return nil
}

func DocumentDeserialize(this js.Value, i []js.Value) interface {} {
	byteArray := i[0]
	initCallback := i[1]

	serialized := make([]byte, byteArray.Get("length").Int())

	js.CopyBytesToGo(serialized, i[0])

	doc := crdt.NewBasicDocument(crdt.NewBasicPositionManager())
	err := doc.Deserialize(serialized)
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
	doc := crdt.NewBasicDocument(crdt.NewBasicPositionManager())

	fd := docManager.PutDocument(doc)

	return int(fd)
}

func DocumentClose(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()

	docManager.RemoveDocument(FileDescriptor(fd))
	return nil
}

func DocumentInsertAt(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()
	value := i[1].String()
	index := i[2].Int()

	doc, err := docManager.GetDocument(FileDescriptor(fd))

	if err != nil {
		fmt.Println(err)
		return -1
	}

	doc.InsertAtIndex(value, index, 1)

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
		return nil
	}

	res := make([]interface{}, len(serialized))

	for i, val := range serialized {
		res[i] = val
	}

	return res
}

func DocumentAddListener(this js.Value, i []js.Value) interface {} {

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
	js.Global().Set("DocumentAddListener", js.FuncOf(DocumentAddListener))
}

func main() {
	c := make(chan struct{}, 0)

	docManager = NewDocumentManager()
	registerCallbacks()
	fmt.Println("Callbacks registered")

	<- c // prevent main from terminating
}