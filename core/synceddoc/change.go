package synceddoc

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

const CHANGE_INSERT      = "insert"
const CHANGE_DELETE      = "delete"
const CHANGE_PEER_CURSOR = "peer_cursor"

type ChangeInsert struct {
	Value string
	Index int
}

type ChangeCRDTInsert struct {
	Value    string
	Position crdt.Position
}

type ChangeDelete struct {
	Index int
}

type ChangeCRDTDelete struct {
	Position crdt.Position
}

type ChangePeerCursor struct {
	PeerID          utils.UUID
	CursorPosition  int
}