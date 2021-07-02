package network

type SyncManager interface {

	SetOnChangeListener()
	SetPeerConnectedListener()
	SetPeerDisconnectedListener()

	BroadcastChange(change interface {})

	// Connect establishes necessary connections and enables
	// receiving and sending changes to and from network.
	// Applications must set listeners using SetListener
	// before calling Start
	Connect()

	// Disconnect terminates established connections and frees resources
	Disconnect()
}