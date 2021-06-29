package main

import (
	"fmt"
	"syscall/js"
)

func DocumentOpen(this js.Value, i []js.Value) interface {} {

	return nil
}

func DocumentDeserialize(this js.Value, i []js.Value) interface {} {

	return nil
}

func DocumentNew(this js.Value, i []js.Value) interface {} {

	return nil
}

func DocumentClose(this js.Value, i []js.Value) interface {} {

	return nil
}

func DocumentInsertAt(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()
	value := i[1].String()
	index := i[2].Int()

	fmt.Printf("Calling DocumentInsertAt on fd: %v index: %d, value: %s\n", fd, index, value)
	return nil
}

func DocumentDeleteAt(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()
	index := i[1].Int()

	fmt.Printf("Calling DocumentDeleteAt on fd: %v index: %d\n", fd, index)
	return nil
}

func DocumentChangeCursor(this js.Value, i []js.Value) interface {} {
	fd := i[0].Int()
	index := i[1].Int()

	fmt.Printf("Calling DocumentChangeCursor on fd: %v index: %d\n", fd, index)
	return nil
}

func DocumentSerialize(this js.Value, i []js.Value) interface {} {

	return nil
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

	fmt.Println("Hello world")
	registerCallbacks()

	<- c // prevent main from terminating
}