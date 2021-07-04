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

func FromBytes(data [][]byte) PackedDocument {
	packedDocument := PackedDocument{}
	packedDocument.Site = string(data[0])
	packedDocument.Position = string(data[1])
	packedDocument.Value = string(data[2])
	packedDocument.Action = string(data[3])
	return packedDocument
}


func GetPackedDocuments(data []byte) []PackedDocument {
	var documents []PackedDocument
	var dataToConvert [][]byte
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		dataToConvert = append(dataToConvert, line)
		if len(dataToConvert) == 4 {
			var packedDocument = FromBytes(dataToConvert)
			documents = append(documents, packedDocument)
			dataToConvert = dataToConvert[:0]
		}
	}
	return documents
}

func GetCopy(packedDocument PackedDocument) PackedDocument {
	copy := PackedDocument{packedDocument.Site, packedDocument.Position, packedDocument.Value, packedDocument.Action}
	return copy
}