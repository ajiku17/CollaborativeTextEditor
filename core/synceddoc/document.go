package synceddoc

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

// ChangeListener
// change is one of several types:
//  ------------------
// | MessageInsert    |
// | MessageDelete    |
type ChangeListener func (changeName string, change interface {}, aux interface{})

type PeerConnectedListener func (peerId utils.UUID, cursorPosition int, aux interface{})
type PeerDisconnectedListener func (peerId utils.UUID, aux interface{})

type Op interface{}

type Document interface {
	GetID() utils.UUID

	// ConnectSignals sets signal handlers
	ConnectSignals(changeListener ChangeListener,
		peerConnectedListener PeerConnectedListener,
		peerDisconnectedListener PeerDisconnectedListener)

	/*
	 * Atomically sets listeners, without missing any changes.
	 * Applications should use these functions when they wish to
	 * change the listeners already passed in at Document creation time.
	 */
	SetChangeListener(listener ChangeListener)
	SetPeerConnectedListener(listener PeerConnectedListener)
	SetPeerDisconnectedListener(listener PeerDisconnectedListener)

	/*
	 * Returns the contents of this document serialized into a byte array
	 */
	Serialize() ([]byte, error)

	/*
	 * Document modifications
	 */

	// LocalInsert currently only supports strings of length 1
	LocalInsert(index int, val string)
	LocalDelete(index int)
	ApplyRemoteOp(peerId utils.UUID, op Op, aux interface{})
	SetCursor(index int)

	/*
	 * Closes the document, frees resources. Document becomes non editable.
	 */
	Close()

	ToString() string
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

