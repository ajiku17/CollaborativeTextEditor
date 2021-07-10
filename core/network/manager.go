package network

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type Manager interface {
	GetId() utils.UUID
	// Start establishes necessary connections and enables
	// receiving and sending changes to and from network.
	// Applications must set listeners using SetListener
	// before calling Start
	Start()

	// Stop terminates established connections
	Stop()

	// Kill frees resources and
	Kill()
}