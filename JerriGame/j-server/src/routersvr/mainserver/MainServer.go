package mainserver

import (
	"fmt"
	"time"

	"encoding/binary"

	message "jserver/src/common/message"
	"jserver/src/protocol/ss"

	gnet "github.com/walkon/wsgnet"

	conf "routersvr/config"
	server "routersvr/server"
)

type MainServer struct {
	//
	eng gnet.Engine
	// 网络类型
	network string
	// 地址
	addr string
	// 单核
	multicore bool

	async bool
	codec gnet.ICodec

	polltimes int

	connIdx uint64
}

// PollerPreInit is called before the poller starts.
func (s *MainServer) PollerPreInit() {
	fmt.Println("PollerPreInit")
}

// PollerPostInit is called after the poller starts.
func (s *MainServer) PollerProc() error {
	// fmt.Println("PollerProc polltimes: ", s.polltimes)

	return nil
}

// PollerWaitTimeOut is called when the poller is waiting for events.
func (s *MainServer) PollerWaitTimeOut() int {
	s.polltimes = s.polltimes + 1
	fmt.Println("PollerWaitTimeOut polltimes: ", s.polltimes)

	return 10
}

func (s *MainServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	fmt.Println("OnBoot")

	s.eng = eng

	return
}

func (s *MainServer) OnShutdown(eng gnet.Engine) {
	fmt.Println("OnShutdown")

}

func (s *MainServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("OnOpen")

	if c.Context() != nil {
		fmt.Println("OnOpen Context: ", c.Context())
		return
	}

	s.connIdx = s.connIdx + 1
	c.SetContext(s.connIdx)
	c.SetWebSock(true)

	fmt.Println("OnOpen connIdx: ", s.connIdx, " Context: ", c.Context(), " RemoteAddr: ", c.RemoteAddr())

	// 创建服务器
	newserver := &server.Server{}
	newserver.SetConnIdx(s.connIdx)
	newserver.SetConn(c)

	// 添加服务器
	server.ServerMgr.AddServerByConnIdx(newserver)

	return
}

func (s *MainServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("OnClose")

	connIdx, ok := c.Context().(uint64)
	if !ok {
		fmt.Println("connIdx Error")
	}

	fmt.Println("OnClose connIdx: ", connIdx, " RemoteAddr: ", c.RemoteAddr())

	// 删除服务器
	server.ServerMgr.DelServerByConnIdx(connIdx)

	return
}

func (s *MainServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	fmt.Println("OnTraffic")

	connIdx, ok := c.Context().(uint64)
	if !ok {
		fmt.Println("React connIdx Error")
	}

	frame, err := c.Peek(-1)
	if err != nil {
		fmt.Println("Peek Error: ", err)
		action = gnet.Close
		return
	}

	fmt.Println("OnTraffic connIdx: ", connIdx, " RemoteAddr: ", c.RemoteAddr(), " frame len: ", len(frame))

	total_len, ssPkg, body, err := message.SSParseProtocolPacket(frame)
	if err != nil {
		fmt.Println("DecodeServerMessage Error: ", err)
		action = gnet.Close
	}
	c.Discard(int(total_len))

	fmt.Println("React connIdx: ", connIdx, " RemoteAddr: ", c.RemoteAddr(), " total len: ", total_len, " cmd: ", ssPkg.Head.Cmd)
	if connIdx == 0 {
		fmt.Println("connIdx Error")
		return
	}
	// 注册
	if ss.SS_CMD(ssPkg.Head.Cmd) == ss.SS_CMD_REGISTER_SERVER_REQ {
		req := body.(*ss.SsCmdRegisterServerReq)
		fmt.Println("REGISTER_SERVER_REQ server_id: ", req.ServerId)

		res_server := server.ServerMgr.GetServerByConnIdx(connIdx)
		if res_server == nil {
			fmt.Println("Server Not Found")
		}
		res_server.SetServerID(uint32(req.ServerId))
		server.ServerMgr.AddServerByServerID(res_server)

		// 返回注册结果
		var res ss.SsCmdRegisterServerRes
		message.SendToRouter(res_server.GetConn(), ss.SS_CMD_REGISTER_SERVER_RES, int32(conf.Conf.ServerInfo.ServerID), req.ServerId, &res)

		return
	} else {
		// 转发服务器消息
		target_server_id := ssPkg.Head.TargetServerId
		target_server := server.ServerMgr.GetServerByServerID(uint32(target_server_id))
		if target_server == nil {
			fmt.Println("Target Server Not Found")
			return
		}

		target_server.SendMsg(frame)
	}

	return
}

func (s *MainServer) OnTick() (delay time.Duration, action gnet.Action) {
	fmt.Println("OnTick")

	return
}

func StartServer() {
	fmt.Println("Start Server")

	reuseport := true
	multicore := false
	async := true
	network := "tcp"
	addr := conf.Conf.NetWork.ClientService

	encoderConfig := gnet.EncoderConfig{
		ByteOrder:         binary.BigEndian,
		LengthFieldLength: 4,
	}

	decoderConfig := gnet.DecoderConfig{
		ByteOrder:         binary.BigEndian,
		LengthFieldLength: 4,
	}

	codec := gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)

	mainserver := &MainServer{
		network:   network,
		addr:      addr,
		multicore: multicore,
		async:     async,
		codec:     codec,
		connIdx:   0,
	}

	fmt.Println("Start Server Success", addr)

	gnet.Run(mainserver, network+"://"+addr,
		gnet.WithLockOSThread(async),
		gnet.WithMulticore(multicore),
		gnet.WithReusePort(reuseport),
		gnet.WithTicker(false),
		gnet.WithTCPKeepAlive(time.Minute*1),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
		gnet.WithLoadBalancing(gnet.SourceAddrHash),
	)
}
