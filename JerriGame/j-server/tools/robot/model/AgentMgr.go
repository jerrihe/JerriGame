package model

import (
	"fmt"
	"sync"
	"time"

	conf "robot/config"

	snow "jserver/src/common/snowflake"

	gnet "github.com/walkon/wsgnet"
)

type AgentMgr struct {
	Agents     map[string]*Agent
	AgentSLock sync.RWMutex

	ConnIdx2Agent map[uint64]*Agent

	Eng     gnet.Engine
	EngLock sync.RWMutex
	IsEng   bool

	is_init bool
}

var AgentMgrIn *AgentMgr

func init() {
	fmt.Printf("AgentMgr init\n")
	AgentMgrIn = &AgentMgr{}
}

func (am *AgentMgr) Init() {
	am.is_init = false
	am.IsEng = false

	am.Agents = make(map[string]*Agent)
	am.ConnIdx2Agent = make(map[uint64]*Agent)

	am.InitAgent()
}

func (am *AgentMgr) InitAgent() {
	for i := int(conf.Conf.ServerInfo.AgentMinNum); i <= int(conf.Conf.ServerInfo.AgentMaxNum); i++ {
		fmt.Println("InitAgent user: ", i, " plat_type: 1")
		conn_id := snow.SF.NextID()
		AgentMgrIn.AddAgent(fmt.Sprint(i), "1", uint64(conn_id), nil)
	}
}

func (am *AgentMgr) Tick1s(now time.Time) {

	// if !am.is_init {
	// 	am.is_init = true
	// 	for i := int(conf.Conf.ServerInfo.AgentMinNum); i <= int(conf.Conf.ServerInfo.AgentMaxNum); i++ {
	// 		conId := snow.SF.NextID()
	// 		conn, cerr := gnet.AddTCPConnector(&am.Eng, "tcp", conf.Conf.ServerInfo.ServerAddr, uint64(conId), gnet.WithTCPNoDelay(gnet.TCPDelay))
	// 		if cerr != nil {
	// 			fmt.Println("AddTCPConnector error")
	// 			continue
	// 		}

	// 		AgentMgrIn.AddAgent(fmt.Sprint(i), "1", uint64(conId), conn)
	// 	}
	// }

	// am.AgentSLock.RLock()
	// defer am.AgentSLock.RUnlock()

	// for _, agent := range am.Agents {
	// 	agent.Tick1s(now)
	// }
}

func (am *AgentMgr) SetEngine(eng gnet.Engine) {
	am.EngLock.Lock()
	defer am.EngLock.Unlock()

	am.Eng = eng
	am.IsEng = true
}

func (am *AgentMgr) GetEngine() gnet.Engine {
	am.EngLock.RLock()
	defer am.EngLock.RUnlock()

	return am.Eng
}

func (am *AgentMgr) AddAgent(user string, plat_type string, connIdx uint64, conn gnet.Conn) {
	am.AgentSLock.Lock()
	defer am.AgentSLock.Unlock()

	agent := &Agent{}
	agent.Init(user, plat_type, connIdx, conn)

	am.Agents[user+"_"+plat_type] = agent
	am.ConnIdx2Agent[connIdx] = agent

	// 增加协程处理
	go Active(agent)
}

func (am *AgentMgr) GetAgent(user string, plat_type string) *Agent {
	am.AgentSLock.RLock()
	defer am.AgentSLock.RUnlock()

	return am.Agents[user+"_"+plat_type]
}

func (am *AgentMgr) GetAgentByConnIdx(connIdx uint64) *Agent {
	am.AgentSLock.RLock()
	defer am.AgentSLock.RUnlock()

	return am.ConnIdx2Agent[connIdx]
}

func (am *AgentMgr) DelAgent(user string, plat_type string, connIdx uint64) {
	am.AgentSLock.Lock()
	defer am.AgentSLock.Unlock()

	delete(am.Agents, user+"_"+plat_type)
	delete(am.ConnIdx2Agent, connIdx)
}

func InitAgent() {
	// for i := 0; i < int(conf.Conf.ServerInfo.AgentNum); i++ {
	// 	conId := snow.SF.NextID()
	// 	conn, cerr := gnet.AddTCPConnector(&am.Eng, "tcp", conf.Conf.ServerInfo.ServerAddr, uint64(conId), gnet.WithTCPNoDelay(gnet.TCPDelay))
	// 	if cerr != nil {
	// 		fmt.Println("AddTCPConnector error")
	// 		continue
	// 	}

	// 	AgentMgrIn.AddAgent(fmt.Sprint(i), "1", uint64(conId), conn)
	// }
}
