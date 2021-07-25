package network

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type PeerConnectedListener func (peerId utils.UUID, cursorPosition int, aux interface{})
type PeerDisconnectedListener func (peerId utils.UUID, aux interface{})

type Manager interface {
	GetId() utils.UUID
	// Start establishes necessary connections and enables
	// receiving and sending changes to and from network.
	// Applications must set listeners using SetListener
	// before calling Start
	Start()

	// Stop terminates established connections
	Stop()

	ConnectSignals(peerConnectedListener PeerConnectedListener,
		peerDisconnectedListener PeerDisconnectedListener)

	OnPeerConnect(peerConnectedListener PeerConnectedListener)
	OnPeerDisconnect(peerDisconnectedListener PeerDisconnectedListener)

	// Kill frees resources and
	Kill()
}