package synceddoc

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

const CHANGE_INSERT      = "insert"
const CHANGE_DELETE      = "delete"
const CHANGE_PEER_CURSOR = "peer_cursor"
const CONNECT = "connect"

type ChangeInsert struct {
	ManagerId utils.UUID
	Value string
	Index int
}

type ChangeCRDTInsert struct {
	ManagerId utils.UUID
	Value    string
	Position crdt.Position
}

type ChangeDelete struct {
	ManagerId utils.UUID
	Index int
}

type ChangeCRDTDelete struct {
	ManagerId utils.UUID
	Position crdt.Position
}

type ChangePeerCursor struct {
	PeerID          utils.UUID
	CursorPosition  int
}

type ConnectRequest struct {
	Id          utils.UUID
}