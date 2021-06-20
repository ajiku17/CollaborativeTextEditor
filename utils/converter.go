package utils

import (
	"bytes"
)

type PackedDocument struct {
	Site string
	Position string
	Value string
	Action string
}

func ToBytes(packedDocument PackedDocument) []byte {
	fullData := packedDocument.Site + "\n" + packedDocument.Position + "\n" + packedDocument.Value + "\n" + packedDocument.Action + "\n"
	return []byte(fullData)
}

func FromBytes(data []byte) PackedDocument {
	tokenizedData := bytes.Split(data, []byte("\n"))
	packedDocument := PackedDocument{}
	packedDocument.Site = string(tokenizedData[0])
	packedDocument.Position = string(tokenizedData[1])
	packedDocument.Value = string(tokenizedData[2])
	packedDocument.Action = string(tokenizedData[3])
	return packedDocument
}