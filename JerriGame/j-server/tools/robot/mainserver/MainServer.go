package mainserver

import (
	"fmt"
	"time"

	"encoding/binary"

	gnet "github.com/walkon/wsgnet"

	"jserver/src/common/message"
	timer "jserver/src/common/time"
	conf "robot/config"
	handle "robot/handle"
	model "robot/model"

	_ "robot/app"
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
	model.AgentMgrIn.SetEngine(eng)

	return
}

func (s *MainServer) OnShutdown(eng gnet.Engine) {
	fmt.Println("OnShutdown")

}

func (s *MainServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("OnOpen")

	if c.Context() != nil {
		connIdx := c.Context().(uint64)
		fmt.Println("OnOpen Context ", connIdx)
		return
	}

	return
}

func (s *MainServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("OnClose")

	connIdx, ok := c.Context().(uint64)
	if !ok {
		fmt.Println("React ConnIdx Error")
	}

	fmt.Println("OnClose ConnIdx: ", connIdx, " RemoteAddr: ", c.RemoteAddr())

	agent := model.AgentMgrIn.GetAgentByConnIdx(connIdx)
	if agent == nil {
		fmt.Println("GetAgentByConnIdx Error")
	}

	agent.SetConnectState(model.ConnectNone)

	return
}

func (s *MainServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	fmt.Println("OnTraffic")

	connIdx, ok := c.Context().(uint64)
	if !ok {
		fmt.Println("React ConnIdx Error", c.Context(), " connIdx: ", connIdx)
		action = gnet.Close
	}

	frame, err := c.Peek(-1)
	if err != nil {
		fmt.Println("Peek Error: ", err)
		action = gnet.Close
		return
	}

	total_len, ssPkg, body, err := message.CSParseProtocolPacket(frame)
	if err != nil {
		fmt.Println("CSParseProtocolPacket Error")
		// action = gnet.Close
		return
	}
	c.Discard(int(total_len))
	player := model.AgentMgrIn.GetAgentByConnIdx(uint64(connIdx))

	if player == nil {
		fmt.Println("GetPlayerByConnIdx Error")
		action = gnet.Close
		return
	}

	handle.HandleMsg(player, ssPkg.Head, body)
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
	model.AgentMgrIn.Tick1s(now)
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
