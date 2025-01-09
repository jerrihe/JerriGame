package model

import (
	"fmt"
	"jserver/src/common/message"
	"jserver/src/protocol/cs"
	"math/rand/v2"
)

/*
	type ActiveMap struct {
		// 节点ID
		NodeID uint32
		// 节点类型
		NodeType uint32
		// 子节点列表
		NodeList []uint32
		// 执行动作
		Action func()
		// 执行状态
		NodeState uint32
	}
*/
var GridMaxRow = 30
var GridMaxCol = 30
var GridHeight = 5

var SceneActiveMap map[uint32]ActiveMap

func init() {
	SceneActiveMap = make(map[uint32]ActiveMap)
}

func InitScene() {
	SceneActiveMap = map[uint32]ActiveMap{
		1: {NodeType: ACTIVE_NODE_TYPE_ORDER, NodeList: []uint32{2, 3}, Action: EnterScene},
		2: {NodeType: ACTIVE_NODE_TYPE_SUB, Action: EnterScene},
		3: {NodeType: ACTIVE_NODE_TYPE_SUB, Action: MoveScene},
		// 4: {NodeType: ACTIVE_NODE_TYPE_SUB, Action: LeaveScene},
	}

	// 注册行为
	RegisterActive(ACTIVE_TYPE_SCENE, SceneActiveMap)
}

func EnterScene(agent *Agent) {
	fmt.Println("app EnterScene")

	var req cs.CsCmdSceneEnterReq
	req.SceneId = 1
	req.MapId = 1

	message.SendToClientMsg(agent.GetConn(), int32(cs.CS_CMD_SCENE_ENTER_REQ), 0, &req)
}

func LeaveScene(agent *Agent) {
	fmt.Println("app LeaveScene")

	var req cs.CsCmdSceneLeaveReq
	req.SceneId = 1
	req.MapId = 1

	message.SendToClientMsg(agent.GetConn(), int32(cs.CS_CMD_SCENE_LEAVE_REQ), 0, &req)
}

func MoveScene(agent *Agent) {
	fmt.Println("app MoveScene")

	var req cs.CsCmdSceneMoveReq
	req.SceneId = agent.SceneId
	req.MapId = agent.MapId
	// 随机移动
	min := 0.0
	x_max := float64(GridHeight * GridMaxRow)
	y_max := float64(GridHeight * GridMaxCol)
	xRandomFloat := min + rand.Float64()*(x_max-min)
	yRandomFloat := min + rand.Float64()*(y_max-min)
	req.PosX = uint32(xRandomFloat * 100)
	req.PosY = uint32(yRandomFloat * 100)

	fmt.Println("MoveScene PosX: ", req.PosX, " PosY: ", req.PosY)

	message.SendToClientMsg(agent.GetConn(), int32(cs.CS_CMD_SCENE_MOVE_REQ), 0, &req)
}
