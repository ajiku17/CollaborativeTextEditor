package server

type DummyServer struct {
	Server
}

// func (server *Server) IsConnected(id utils.UUID) bool {
// 	for curr_id, _ := range server.ConnectedSockets {
// 		if curr_id == id {
// 			return true
// 		}
// 	}
// 	return false
// }