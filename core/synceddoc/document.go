package synceddoc

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

// OnChangeListener
// change is one of several types:
//  ------------------
// | ChangeInsert     |
// | ChangeDelete     |
// | ChangePeerCursor |
type OnChangeListener func (changeName string, change interface {})

type PeerConnectedListener func (peerId utils.UUID, cursorPosition int)
type PeerDisconnectedListener func (changeName string, change interface {})

type Document interface {
	GetID() utils.UUID

	/*
	 * Sets listeners
	 */
	SetOnChangeListener(listener OnChangeListener)
	SetPeerConnectedListener(listener PeerConnectedListener)
	SetPeerDisconnectedListener(listener PeerDisconnectedListener)

	/*
	 * Returns the contents of this document serialized into a byte array
	 */
	Serialize() []byte

	/*
	 * Document modifications
	 */
	InsertAtIndex(index int)
	DeleteAtIndex(index int)
	SetCursor(index int)

	/*
	 * Closes the document
	 */
	Close()
}

/*
	concrete implementation example

	type DocumentImplementation struct {
		id           utils.UUID

		document     crdt.Document
		syncManager  network.DocumentSyncManager

		// ... more declarations
	}
*/

