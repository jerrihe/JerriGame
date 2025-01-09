package server

import (
	"fmt"
)

type ServerInfoMgr struct {
	ConnIdx2Server  map[uint64]*Server
	serverID2Server map[uint32]*Server
}

var ServerMgr *ServerInfoMgr

func init() {
	fmt.Println("ServerMgr Init")
	ServerMgr = &ServerInfoMgr{}
	ServerMgr.serverID2Server = make(map[uint32]*Server)
	ServerMgr.ConnIdx2Server = make(map[uint64]*Server)
}

func (s *ServerInfoMgr) AddServerByConnIdx(server *Server) {
	s.ConnIdx2Server[server.GetConnIdx()] = server
}

func (s *ServerInfoMgr) AddServerByServerID(server *Server) {
	s.serverID2Server[server.ServerID] = server
}

func (s *ServerInfoMgr) GetServer(serverID uint32) *Server {
	return s.serverID2Server[serverID]
}

func (s *ServerInfoMgr) DelServer(serverID uint32) {
	server := s.serverID2Server[serverID]
	if server != nil {
		delete(s.serverID2Server, serverID)
		delete(s.ConnIdx2Server, server.GetConnIdx())
	}
}

func (s *ServerInfoMgr) GetServerByConnIdx(sessionID uint64) *Server {
	return s.ConnIdx2Server[sessionID]
}

func (s *ServerInfoMgr) DelServerByConnIdx(sessionID uint64) {
	server := s.ConnIdx2Server[sessionID]
	if server != nil {
		delete(s.serverID2Server, server.ServerID)
		delete(s.ConnIdx2Server, sessionID)
	}
}

func (s *ServerInfoMgr) GetServerByServerID(serverID uint32) *Server {
	return s.serverID2Server[serverID]
}
