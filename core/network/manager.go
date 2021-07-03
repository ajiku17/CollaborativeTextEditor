package network

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type DataReceiveListener func (data []byte)
type PeerConnectedListener func (peerId utils.UUID, cursorPosition int)
type PeerDisconnectedListener func (peerId utils.UUID)

type Manager interface {

	SetOnDataReceiveListener(listener DataReceiveListener)
	SetPeerConnectedListener(listener PeerConnectedListener)
	SetPeerDisconnectedListener(listener PeerDisconnectedListener)

	BroadcastChange(change interface {})

	// Connect establishes necessary connections and enables
	// receiving and sending changes to and from network.
	// Applications must set listeners using SetListener
	// before calling Start
	Connect()

	// Disconnect terminates established connections and frees resources
	Disconnect()
}