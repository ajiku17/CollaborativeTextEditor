package main

import (
	"errors"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
)

type FileDescriptor int64

type DocumentManager struct {
	nextFd FileDescriptor
	openDocuments map[FileDescriptor] synceddoc.Document
}

func NewDocumentManager() *DocumentManager {
	manager := new(DocumentManager)

	manager.nextFd = 1
	manager.openDocuments = make(map[FileDescriptor] synceddoc.Document)

	return manager
}


func (manager *DocumentManager) GetNextFd() FileDescriptor {
	res := manager.nextFd
	manager.nextFd++

	return res
}

func (manager *DocumentManager) PutDocument(doc synceddoc.Document) FileDescriptor {
	res := manager.GetNextFd()

	manager.openDocuments[res] = doc

	return res
}

func (manager *DocumentManager) GetDocument(fd FileDescriptor) (synceddoc.Document, error) {
	doc, ok := manager.openDocuments[fd]

	if !ok {
		return nil, errors.New("document not found")
	}

	return doc, nil
}

func (manager *DocumentManager) RemoveDocument(fd FileDescriptor) {
	delete(manager.openDocuments, fd)
}


