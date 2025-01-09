package mainserver

import (
	"fmt"
	"gamesvr/game"
	"time"

	"encoding/binary"

	message "jserver/src/common/message"

	gnet "github.com/walkon/wsgnet"

	conf "gamesvr/config"
	timer "jserver/src/common/time"

	serverMgr "gamesvr/server"

	model "gamesvr/model"
)

type MainServer struct {
	//
	gnet.EventHandler
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

	timer.ProcessTimer()
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

	serverMgr.ServerMgr.InitServer(eng)

	return
}

func (s *MainServer) OnShutdown(eng gnet.Engine) {
	fmt.Println("OnShutdown")
	model.PlayerMgr.KickAllPlayer(1)
}

func (s *MainServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("OnOpen")

	if c.Context() != nil {
		fmt.Println("OnOpen ConnIdx: ", c.Context().(uint64), " Context: ", c.Context(), " RemoteAddr: ", c.RemoteAddr(), " LocalAddr: ", c.LocalAddr())

		server := serverMgr.ServerMgr.GetServerByConnIdx(c.Context().(uint64))
		if server == nil {
			fmt.Println("GetServerByConnIdx Error")
			action = gnet.Close
		}

		server.SetState(1)
		fmt.Println("Connect Success")
		return
	}

	player := model.PlayerMgr.CreatePlayer(c)
	if player == nil {
		fmt.Println("Create Player Error")
		action = gnet.Close
	}

	c.SetContext(uint64(player.GetConnIdx()))
	c.SetWebSock(true)

	fmt.Println("OnOpen conId: ", player.GetConnIdx(), " RemoteAddr: ", c.RemoteAddr())

	return
}

func (s *MainServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("OnClose")

	connIdx, ok := c.Context().(uint64)
	if !ok {
		fmt.Println("React ConnIdx Error")
	}

	fmt.Println("OnClose ConnIdx: ", connIdx, " RemoteAddr: ", c.RemoteAddr())
	if connIdx < 10000 {
		// 服务器内部消息
		serverMgr.ServerMgr.DelServerByConnIdx(connIdx)
	} else {
		player := model.PlayerMgr.GetPlayerByConnIdx(int64(connIdx))
		if player == nil {
			fmt.Println("GetPlayerByConnIdx Error")
		}
		model.PlayerMgr.LoginOut(player, 1)
	}

	return
}

func (s *MainServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	fmt.Println("OnTraffic")

	connIdx, ok := c.Context().(uint64)
	if !ok {
		fmt.Println("React ConnIdx Error", c.Context(), " connIdx: ", connIdx)
		action = gnet.Close
	}

	if connIdx == 0 {
		fmt.Println("ConnIdx Error")
		return
	}

	frame, err := c.Peek(-1)
	if err != nil {
		fmt.Println("Peek Error: ", err)
		action = gnet.Close
		return
	}

	if connIdx < 10000 {
		// 服务器内部消息
		total_len, ssPkg, body, err := message.SSParseProtocolPacket(frame)
		if err != nil {
			fmt.Println("DecodeServerMessage Error: ", err)
			action = gnet.Close
		}
		c.Discard(int(total_len))
		fmt.Println("React ConnIdx: ", connIdx, " RemoteAddr: ", c.RemoteAddr(), " len: ", total_len, " cmd: ", ssPkg.Head.Cmd)

		game.HandleSS(ssPkg.Head, &body)
		return
	}

	// 解包
	total_len, csPkg, body, err := message.CSParseProtocolPacket(frame)
	if err != nil {
		fmt.Println("CSParseProtocolPacket Error")
		// action = gnet.Close
		return
	}

	c.Discard(int(total_len))
	player := model.PlayerMgr.GetPlayerByConnIdx(int64(connIdx))
	if player == nil {
		fmt.Println("GetPlayerByConnIdx Error", connIdx)
		action = gnet.Close
		return
	}

	game.HandleCs(player, csPkg.Head, body)

	return
}

func (s *MainServer) OnTick() (delay time.Duration, action gnet.Action) {
	fmt.Println("OnTick")

	return
}

func (s *MainServer) BuildTicker() {
	// timer.AddTimer("MainServer::tick200ms", 200*time.Millisecond, s.tick200ms)
	timer.AddTimer("MainServer::tick1s", time.Second, s.tick1s)
	// timer.AddTimer("MainServer::tick5s", 5*time.Second, s.tick5s)
}

func (s *MainServer) tick1s(now time.Time) {
	// fmt.Println("MainServer tick1s")
	model.PlayerMgr.Tick1S()
	serverMgr.ServerMgr.Tick1S()
}

func StartServer() {
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

	mainserver.BuildTicker()

	gnet.Run(mainserver, network+"://"+addr,
		gnet.WithLockOSThread(async),
		gnet.WithMulticore(multicore),
		gnet.WithReusePort(reuseport),
		gnet.WithTicker(false),
		gnet.WithTCPKeepAlive(time.Minute*1),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
		gnet.WithLoadBalancing(gnet.SourceAddrHash),
		gnet.WithReusePort(true),
	)
}
