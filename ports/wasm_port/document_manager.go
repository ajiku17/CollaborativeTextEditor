package main

import (
	"errors"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type DocumentID utils.UUID

type DocumentManager struct {
	openDocuments map[DocumentID] synceddoc.Document
}

func NewDocumentManager() *DocumentManager {
	manager := new(DocumentManager)

	manager.openDocuments = make(map[DocumentID] synceddoc.Document)

	return manager
}

func (manager *DocumentManager) PutDocument(doc synceddoc.Document) {
	manager.openDocuments[DocumentID(doc.GetID())] = doc
}

func (manager *DocumentManager) GetDocument(docId DocumentID) (synceddoc.Document, error) {
	doc, ok := manager.openDocuments[docId]

	if !ok {
		return nil, errors.New("document not found")
	}

	return doc, nil
}

func (manager *DocumentManager) RemoveDocument(docId DocumentID) {
	delete(manager.openDocuments, docId)
}


