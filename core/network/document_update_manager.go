package network

import "github.com/ajiku17/CollaborativeTextEditor/utils"

type DocumentUpdateManager interface {
	ConnectWithServer(int)
	Insert(position Position, val string, site int)
	Delete(position Position, site int)
	AddListener()
	Notify() *utils.PackedDocument
}