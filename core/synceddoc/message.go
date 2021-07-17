package synceddoc

import (
	"github.com/ajiku17/CollaborativeTextEditor/core/crdt"
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

const MESSAGE_INSERT      = "insert"
const MESSAGE_DELETE      = "delete"
const MESSAGE_PEER_CURSOR = "peer_cursor"
const MESSAGE_CONNECT     = "connect"

type MessageInsert struct {
	Value string
	Index int
}

type MessageDelete struct {
	Index int
}

type MessageCRDTInsert struct {
	ManagerId utils.UUID
	Value    string
	Position crdt.Position
}

type MessageCRDTDelete struct {
	ManagerId utils.UUID
	Position crdt.Position
}

type MessagePeerCursor struct {
	PeerID          utils.UUID
	CursorPosition  int
}

type ConnectRequest struct {
	Id          utils.UUID
}

type OperationRequest struct {
	Id utils.UUID
	Operation interface{}
}