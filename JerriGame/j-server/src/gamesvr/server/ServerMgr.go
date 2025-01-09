package server

import (
	"fmt"
	"jserver/src/common/message"
	"jserver/src/protocol/ss"

	gnet "github.com/walkon/wsgnet"

	conf "gamesvr/config"
)

type ServerInfoMgr struct {
	ServerMap map[uint64]*Server
	Server    *Server
	Eng       gnet.Engine

	ServerConnId uint64
}

var ServerMgr *ServerInfoMgr

func init() {
	fmt.Printf("ServerMgr init\n")
	ServerMgr = &ServerInfoMgr{}
	ServerMgr.ServerMap = make(map[uint64]*Server)
	ServerMgr.Server = new(Server)
	ServerMgr.ServerConnId = 0

}

func (s *ServerInfoMgr) InitServer(eng gnet.Engine) {
	s.Server.ServerID = conf.Conf.ServerInfo.ServerID
	fmt.Println("ServerID:", s.Server.ServerID, conf.Conf.ServerInfo.ServerID)
	s.Eng = eng

	for _, router_server := range conf.GetRelationService()["RouterSvr"] {
		ServerMgr.ServerConnId++
		new_server := &Server{}
		new_server.State = 0
		new_server.Addr = router_server
		new_server.ConnIdx = s.ServerConnId
		s.ServerMap[s.ServerConnId] = new_server
	}
}

func (s *ServerInfoMgr) Connect() {
	for _, server := range s.ServerMap {
		if server.GetState() == 0 {
			// fmt.Println("Connect Start")
			conn, cerr := gnet.AddTCPConnector(&s.Eng, "tcp", server.Addr, server.ConnIdx, gnet.WithTCPNoDelay(gnet.TCPDelay))
			if cerr != nil {
				// fmt.Println("AddTCPConnector error")
				continue
			}
			server.SetConn(conn)
			// server.SetState(1)
			// fmt.Println("Connect Success")
		}
	}
}

func (s *ServerInfoMgr) Register() {
	for _, server := range s.ServerMap {
		if server.GetState() == 1 {
			// 注册
			fmt.Println("Register")

			// 发送注册消息
			var registerReq ss.SsCmdRegisterServerReq
			registerReq.ServerId = conf.Conf.ServerInfo.ServerID

			message.SendToRouter(server.GetConn(), ss.SS_CMD_REGISTER_SERVER_REQ, conf.Conf.ServerInfo.ServerID, 0, &registerReq)
			server.SetState(2)
		}
	}
}

func (s *ServerInfoMgr) Tick1S() {
	s.Connect()
	s.Register()
}

func (s *ServerInfoMgr) GetServerConn() gnet.Conn {
	for _, server := range s.ServerMap {
		if server.GetState() == 2 {
			return server.GetConn()
		}
	}
	return nil
}

func (s *ServerInfoMgr) GetServerByConnIdx(Idx uint64) *Server {
	return s.ServerMap[Idx]
}

func (s *ServerInfoMgr) DelServerByConnIdx(connIdx uint64) {
	delete(s.ServerMap, connIdx)
}
