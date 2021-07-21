package signaling

const CONN_OFFER    = "offer"
const CONN_ANSWER   = "answer"
const ICE_CANDIDATE = "ice-candidate"

const MESSAGE_SUBSCRIBE = "subscribe"
const MESSAGE_FORWARD = "forward"

type ConnOffer struct {
	Sender string `json:"sender"`
	SDP    string `json:"sdp"`
}

type ConnAnswer struct {
	Sender string `json:"sender"`
	SDP    string `json:"sdp"`
}

type ICECandidate struct {
	Sender       string `json:"sender"`
	IceCandidate string `json:"ice_candidate"`
}

type Subscription struct {
	PeerId string `json:"peer_id"`
	DocId  string `json:"doc_id"`
}

type SignalMessage struct {
	Msg      []byte `json:"msg"`
	MsgType  string `json:"msg_type"`
	Receiver string `json:"receiver"`
}
