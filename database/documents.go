package database

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type DocumentDatabase map[utils.UUID]interface{}

var documentDatabase DocumentDatabase

func GetDatabase() DocumentDatabase {
	if documentDatabase == nil {
		documentDatabase = make(map[utils.UUID]interface{})
	}
	return documentDatabase
}

func (db *DocumentDatabase)GetDocument(id utils.UUID) interface{} {
	if document, ok := (*db)[id]; ok {
		return document
	}
	return nil
}

func (db *DocumentDatabase)AddDocument(id utils.UUID, data interface{}) {
	(*db)[id] = data
}