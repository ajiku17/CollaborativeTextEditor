package utils

import (
	"bytes"
)

type PackedDocument struct {
	Site string
	Index string
	Position string
	Value string
	Action string
}

func ToBytes(packedDocument PackedDocument) []byte {
	fullData := packedDocument.Site + "\n" + packedDocument.Index + "\n" + packedDocument.Position + "\n" + packedDocument.Value + "\n" + packedDocument.Action + "\n"
	return []byte(fullData)
}

func FromBytes(data [][]byte) PackedDocument {
	packedDocument := PackedDocument{}
	packedDocument.Site = string(data[0])
	packedDocument.Index = string(data[1])
	packedDocument.Position = string(data[2])
	packedDocument.Value = string(data[3])
	packedDocument.Action = string(data[4])
	return packedDocument
}


func GetPackedDocuments(data []byte) []PackedDocument {
	var documents []PackedDocument
	var dataToConvert [][]byte
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		dataToConvert = append(dataToConvert, line)
		if len(dataToConvert) == 5 {
			var packedDocument = FromBytes(dataToConvert)
			documents = append(documents, packedDocument)
			dataToConvert = dataToConvert[:0]
		}
	}
	return documents
}