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

	return nil
}

func DocumentDeleteAt(this js.Value, i []js.Value) interface {} {

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
	js.Global().Set("DocumentSerialize", js.FuncOf(DocumentSerialize))
	js.Global().Set("DocumentAddListener", js.FuncOf(DocumentAddListener))
}

func main() {
	fmt.Println("Hello world")
	registerCallbacks()
}