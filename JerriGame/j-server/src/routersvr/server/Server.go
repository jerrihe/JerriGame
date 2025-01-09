package server

import (
	"fmt"

	gnet "github.com/walkon/wsgnet"
)

type Server struct {
	//
	conn       gnet.Conn
	session_id uint64
	ServerID   uint32
}

func (s *Server) SetConnIdx(session_id uint64) {
	s.session_id = session_id
}

func (s *Server) GetConnIdx() uint64 {
	return s.session_id
}

func (s *Server) SetConn(conn gnet.Conn) {
	s.conn = conn
}

func (s *Server) GetConn() gnet.Conn {
	return s.conn
}

func (s *Server) SetServerID(serverID uint32) {
	s.ServerID = serverID
}

func (s *Server) GetServerID() uint32 {
	return s.ServerID
}

func (s *Server) SendMsg(msg []byte) {
	s.conn.Write(msg)
	fmt.Println("SendMsg ServerID:", s.ServerID, "len:", len(msg))
}
