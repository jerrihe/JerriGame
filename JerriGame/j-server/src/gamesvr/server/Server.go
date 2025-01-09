package server

import (
	"fmt"

	gnet "github.com/walkon/wsgnet"
)

type Server struct {
	//
	conn       gnet.Conn
	session_id uint64
	ServerID   int32
	State      int32
	Addr       string
	ConnIdx    uint64
}

func (s *Server) SetSessionID(session_id uint64) {
	s.session_id = session_id
}

func (s *Server) GetSessionID() uint64 {
	return s.session_id
}

func (s *Server) SetConn(conn gnet.Conn) {
	s.conn = conn
}

func (s *Server) GetConn() gnet.Conn {
	return s.conn
}

func (s *Server) SetServerID(serverID int32) {
	s.ServerID = serverID
}

func (s *Server) GetServerID() int32 {
	return s.ServerID
}

func (s *Server) SetState(state int32) {
	s.State = state
}

func (s *Server) GetState() int32 {
	return s.State
}

func (s *Server) SendMsg(msg []byte) {

	fmt.Println("SendMsg", len(msg))
}
