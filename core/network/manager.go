package network

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type MessageReceiveListener func (message interface{})
type PeerConnectedListener func (peerId utils.UUID, aux interface{})
type PeerDisconnectedListener func (peerId utils.UUID, aux interface{})


type Manager interface {

	SetOnMessageReceiveListener(listener MessageReceiveListener)
	SetPeerConnectedListener(listener PeerConnectedListener)
	SetPeerDisconnectedListener(listener PeerDisconnectedListener)

	BroadcastMessage(message interface{})

	// Connect establishes necessary connections and enables
	// receiving and sending changes to and from network.
	// Applications must set listeners using SetListener
	// before calling Start
	Connect()

	// Disconnect terminates established connections
	Disconnect()

	// Kill frees resources and
	Kill()
}