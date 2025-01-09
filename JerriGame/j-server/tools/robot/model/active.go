package model

import (
	"fmt"
	"jserver/src/common/message"
	conf "robot/config"
	"time"

	"jserver/src/protocol/cs"

	gnet "github.com/walkon/wsgnet"
)

var ActiveTypeMap map[uint32]map[uint32]ActiveMap

func init() {
	fmt.Println("Active init")
	// ActiveTypeMap = make(map[uint32]map[uint32]ActiveMap)
}

func InitActive() {
	ActiveTypeMap = make(map[uint32]map[uint32]ActiveMap)
}

func RegisterActive(activeType uint32, active map[uint32]ActiveMap) {
	fmt.Println("RegisterActive", activeType, " active:", active)
	if _, ok := ActiveTypeMap[activeType]; !ok {
		ActiveTypeMap[activeType] = active
	}
}

func RunActiveMap(activeType uint32, nodeId uint32, subNodeIndex uint32, agent *Agent) (uint32, uint32, uint32, uint32) {

	if _, ok := ActiveTypeMap[activeType]; !ok {
		fmt.Println("ActiveTypeMap not have activeType", activeType)
		return activeType, nodeId, subNodeIndex, 0
	}

	activeTypeMap := ActiveTypeMap[activeType]

	if _, ok := activeTypeMap[nodeId]; !ok {
		fmt.Println("ActiveTypeMap not have nodeId", nodeId)
		return activeType, nodeId, subNodeIndex, 0
	}

	node := activeTypeMap[nodeId]
	if node.NodeType == ACTIVE_NODE_TYPE_SUB {
		return activeType, nodeId, subNodeIndex, 0
	}
	if subNodeIndex >= uint32(len(node.NodeList)) {
		fmt.Println("ActiveTypeMap not have subNodeIndex", subNodeIndex)
		return activeType, nodeId, subNodeIndex, 0
	}

	subNodeId := node.NodeList[subNodeIndex]
	if _, ok := activeTypeMap[subNodeId]; !ok {
		fmt.Println("ActiveTypeMap not have subNode", subNodeId)
		return activeType, nodeId, subNodeIndex, 0
	}

	subNode := activeTypeMap[subNodeId]
	if subNode.NodeType != ACTIVE_NODE_TYPE_SUB {
		return activeType, nodeId, subNodeIndex, 0
	}

	subNode.Action(agent)

	// 寻找下一个节点执行点
	return activeType, nodeId, subNodeIndex, subNodeId
}

func SelectNextNode(activeType uint32, nodeId uint32, subNodeIndex uint32, subNodeId uint32) (uint32, uint32, uint32) {
	nextSubNodeIndex := subNodeIndex + 1
	if nextSubNodeIndex >= uint32(len(ActiveTypeMap[activeType][nodeId].NodeList)) {
		nextSubNodeIndex = 0
	} else {
		// 找到了返回
		return activeType, nodeId, nextSubNodeIndex
	}

	// 所有子节点执行完毕，查看是否有下一个节点
	nextNodeId := subNodeId + 1
	if _, ok := ActiveTypeMap[activeType][nextNodeId]; !ok {
		nextNodeId = 1
	} else {
		// 找到了返回
		return activeType, nextNodeId, nextSubNodeIndex
	}

	// 找下个模块类型
	nextActiveType := activeType + 1
	if _, ok := ActiveTypeMap[nextActiveType]; !ok {
		nextActiveType = 1
	} else {
		// 找到了返回
	}
	return nextActiveType, nextNodeId, nextSubNodeIndex
}

func Active(agent *Agent) {
	fmt.Println("go app Active", agent.User)
	for {
		switch agent.ConnectState {
		case ConnectNone:
			// fmt.Println("Agent State 0")
			// 尝试链接
			ConnectServer(agent)
		case ConnectLogin:
			return
		case ConnectGame:
			// 登录成功
			// a.State = ActiveEnumLoginSuccess
			// 链接成功 进行游戏操作
			switch agent.State {
			case ActiveEnumLogin:
				// 开始登录
				fmt.Printf("开始登录 %s\n", agent.User)
				Login(agent)
			case ActiveEnumLoginSuccess:
				// 登录成功 执行操作
				// 下线
				fmt.Printf("登录成功 开始处理业务节点 %s\n", agent.User)
				// LoginOut(a)
				RunActive(agent)
			case ActiveEnumCreateRole:
				// 创建角色
				fmt.Printf("开始创建角色 %s\n", agent.User)
				CreateRole(agent)
			case ActiveEnumCreateRoleSuccess:
				// 创建角色
				fmt.Printf("创建角色成功 开始登录 %s\n", agent.User)
				Login(agent)
			case ActiveEnumLoginOutSuccess:
				// 登出成功
				// 重新登录
				fmt.Printf("登出成功 重新登录 %s\n", agent.User)
				Login(agent)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

}

func RunActive(agent *Agent) {
	fmt.Println("RunActive before", agent.User, " ActiveType: ", agent.ActiveType, " ActiveNodeID: ", agent.ActiveNodeID, " ActiveSubNodeIndex: ", agent.ActiveSubNodeIndex)
	activeType, nodeId, subNodeIndex, subNodeId := RunActiveMap(agent.ActiveType, agent.ActiveNodeID, agent.ActiveSubNodeIndex, agent)
	agent.ActiveType, agent.ActiveNodeID, agent.ActiveSubNodeIndex = SelectNextNode(activeType, nodeId, subNodeIndex, subNodeId)
	fmt.Println("RunActive after", agent.User, " ActiveType: ", agent.ActiveType, " ActiveNodeID: ", agent.ActiveNodeID, " ActiveSubNodeIndex: ", agent.ActiveSubNodeIndex)
}

func Login(agent *Agent) {
	fmt.Println("app login")

	var req cs.CsCmdLoginReq

	req.User = agent.GetUser()
	req.Platform = agent.GetPlatform()
	// 正在登录状态
	agent.SetState(ActiveEnumLoginIng)

	message.SendToClientMsg(agent.GetConn(), int32(cs.CS_CMD_LOGIN_REQ), 0, &req)
}

func CreateRole(agent *Agent) {
	fmt.Println("app CreateRole")

	var req cs.CsCmdCreateAccountReq
	req.User = agent.GetUser()
	req.Platform = agent.GetPlatform()

	// 正在创建角色状态
	agent.SetState(ActiveEnumCreateRoleIng)

	message.SendToClientMsg(agent.GetConn(), int32(cs.CS_CMD_CREATE_ACCOUNT_REQ), 0, &req)
}

func LoginOut(agent *Agent) {
	fmt.Println("app LoginOut")

	var req cs.CsCmdLoginOutReq
	message.SendToClientMsg(agent.GetConn(), int32(cs.CS_CMD_LOGIN_OUT_REQ), 0, &req)
	agent.SetState(ActiveEnumLoginOutIng)
}

func ConnectServer(agent *Agent) {
	fmt.Println("app ConnectServer")

	if !AgentMgrIn.IsEng {
		fmt.Println("AgentMgrIn not have Eng")
		return
	}

	eng := AgentMgrIn.GetEngine()

	conn, cerr := gnet.AddTCPConnector(&eng, "tcp", conf.Conf.ServerInfo.ServerAddr, uint64(agent.ConnIdx), gnet.WithTCPNoDelay(gnet.TCPDelay))
	if cerr != nil {
		fmt.Println("AddTCPConnector error")
		return
	}

	agent.SetConn(conn)
	agent.SetConnectState(ConnectGame)
	agent.SetState(ActiveEnumLogin)

	fmt.Printf("ConnectServer %s ConnectState %d State %d\n", agent.User, agent.ConnectState, agent.State)
}
