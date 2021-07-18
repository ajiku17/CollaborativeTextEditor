package main

import (
	"errors"
	"github.com/ajiku17/CollaborativeTextEditor/core/network"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type DocumentID utils.UUID

type Document struct {
	Doc   synceddoc.Document
	NetManager  network.Manager
}

type DocumentManager struct {
	openDocuments map[DocumentID] Document
}

func NewDocumentManager() *DocumentManager {
	manager := new(DocumentManager)

	manager.openDocuments = make(map[DocumentID] Document)

	return manager
}

func (manager *DocumentManager) PutDocument(doc Document) {
	for _, d := range manager.openDocuments {
		d.NetManager.Stop()
		d.Doc.Close()
	}

	manager.openDocuments[DocumentID(doc.Doc.GetID())] = doc
}

func (manager *DocumentManager) GetDocument(docId DocumentID) (Document, error) {
	doc, ok := manager.openDocuments[docId]

	if !ok {
		return Document{}, errors.New("document not found")
	}

	return doc, nil
}

func (manager *DocumentManager) RemoveDocument(docId DocumentID) {
	d, ok := manager.openDocuments[docId]
	if ok {
		d.NetManager.Stop()
		d.Doc.Close()
	}
	delete(manager.openDocuments, docId)
}


