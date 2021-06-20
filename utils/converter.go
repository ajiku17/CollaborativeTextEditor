package utils

import (
	"bytes"
	"strconv"
)

type PackedDocument struct {
	Site string
	Position string
	Value string
	Action string
}

func ToBytes(site int, position string, val string, action string) []byte {
	fullData := strconv.Itoa(site) + "\n" + position + "\n" + val + "\n" + action + "\n"
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