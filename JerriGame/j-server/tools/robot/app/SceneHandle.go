package app

import (
	"fmt"
	"jserver/src/protocol/cs"
	handle "robot/handle"
	model "robot/model"

	"google.golang.org/protobuf/proto"
)

func init() {
	fmt.Println("app init")

	handle.RegisterHandleFunc(cs.CS_CMD_SCENE_ENTER_RES, HandleCsCmdSceneEnterRes)
	handle.RegisterHandleFunc(cs.CS_CMD_SCENE_LEAVE_RES, HandleCsCmdSceneLeaveRes)
	handle.RegisterHandleFunc(cs.CS_CMD_SCENE_MOVE_RES, HandleCsCmdSceneMoveRes)
}

func HandleCsCmdSceneEnterRes(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	fmt.Println("HandleCsCmdSceneEnterRes")

	var res = (*msg).(*cs.CsCmdSceneEnterRes)
	if res == nil {
		fmt.Println("HandleCsCmdSceneEnterRes res is nil")
		return
	}

	agent.SetSceneInfo(res.MapId, res.SceneId, float32(res.PosX)/100, float32(res.PosY)/100)
}

func HandleCsCmdSceneLeaveRes(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	fmt.Println("HandleCsCmdSceneLeaveRes")

	var res = (*msg).(*cs.CsCmdSceneLeaveRes)
	if res == nil {
		fmt.Println("HandleCsCmdSceneLeaveRes res is nil")
		return
	}

	agent.SetSceneInfo(0, 0, 0, 0)
}

func HandleCsCmdSceneMoveRes(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	fmt.Println("HandleCsCmdSceneMoveRes")

	var res = (*msg).(*cs.CsCmdSceneMoveRes)
	if res == nil {
		fmt.Println("HandleCsCmdSceneMoveRes res is nil")
		return
	}

	agent.SetSceneInfo(res.MapId, res.SceneId, float32(res.PosX)/100, float32(res.PosY)/100)
}
