package network

import (
	"bytes"
	"encoding/gob"
	"github.com/ajiku17/CollaborativeTextEditor/core/synceddoc"
)

type P2PMessage struct {
	Sender   string
	Receiver string
	IsPatch  bool
	Patch    synceddoc.Patch
}

func EncodeP2PMessage(msg P2PMessage) ([]byte, error) {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)

	err := e.Encode(msg)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func DecodeP2PMessage(msg []byte) (P2PMessage, error) {
	r := bytes.NewBuffer(msg)
	d := gob.NewDecoder(r)

	res := P2PMessage{}
	err := d.Decode(&res)
	if err != nil {
		return P2PMessage{}, err
	}

	return res, nil
}
