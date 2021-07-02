package synceddoc

import "github.com/ajiku17/CollaborativeTextEditor/utils"

const CHANGE_INSERT      = "insert"
const CHANGE_DELETE      = "delete"
const CHANGE_PEER_CURSOR = "peer_cursor"

type ChangeInsert struct {
	Value string
	Index int
}

type ChangeDelete struct {
	Index int
}

type ChangePeerCursor struct {
	PeerID          utils.UUID
	CursorPosition  int
}