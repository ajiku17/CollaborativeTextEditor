package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func ToBytes(message interface{}) []byte {
	var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(&message)
    if err != nil {
		fmt.Println(err)
        return nil
    }
    return buf.Bytes()
}


func FromBytes(bytesToConvert []byte) interface{} {
	buf := bytes.NewBuffer(bytesToConvert)

    dec := gob.NewDecoder(buf)

    var i interface{}

    if err := dec.Decode(&i); err != nil {
		fmt.Println(err)
		return nil
    }
	return i
}