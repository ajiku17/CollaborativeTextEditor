package signaling

const MESSAGE_SUBSCRIBE = "subscribe"
const MESSAGE_FORWARD = "forward"

type Subscription struct {
	PeerId string `json:"peer_id"`
	DocId  string `json:"doc_id"`
}

type SignalMessage struct {
	Msg      []byte `json:"msg"`
	MsgType  string `json:"msg_type"`
	Receiver string `json:"receiver"`
}
