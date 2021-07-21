package p2p

const CONN_OFFER    = "offer"
const CONN_ANSWER   = "answer"
const ICE_CANDIDATE = "ice-candidate"

type ConnOffer struct {
	Sender string `json:"sender"`
	SDP    []byte `json:"sdp"`
}

type ConnAnswer struct {
	Sender string `json:"sender"`
	SDP    []byte `json:"sdp"`
}

type ICECandidateMsg struct {
	Sender       string `json:"sender"`
	IceCandidate string `json:"ice_candidate"`
}

type P2PMessage struct {
	MsgType string `json:"msg_type"`
	Msg     []byte `json:"Msg"`
}

