package app

import (
	"fmt"
	handle "gamesvr/game"
	model "gamesvr/model"
	scene "gamesvr/scene"
	"jserver/src/protocol/cs"

	"google.golang.org/protobuf/proto"
)

func init() {
	fmt.Println("Scene Handle init")
	handle.RegisterHandle(cs.CS_CMD_SCENE_ENTER_REQ, HandleCsCmdSceneEnterReq)
	handle.RegisterHandle(cs.CS_CMD_SCENE_LEAVE_REQ, HandleCsCmdSceneLeaveReq)
	handle.RegisterHandle(cs.CS_CMD_SCENE_MOVE_REQ, HandleCsCmdSceneMoveReq)

}

func HandleCsCmdSceneEnterReq(player *model.Player, head *cs.CsHead, body *proto.Message) {
	if player == nil || head == nil || body == nil {
		return
	}

	fmt.Println("HandleCsCmdSceneEnterReq")

	req := (*body).(*cs.CsCmdSceneEnterReq)
	if req == nil {
		fmt.Println("HandleCsCmdSceneEnterReq req is nil")
		return
	}

	map_id := req.GetMapId()
	scene_id := req.GetSceneId()

	if !scene.SceneMgrIn.EnterPlayer(player, map_id, scene_id) {
		player.SendErrCode(cs.CS_CMD_SCENE_ENTER_RES, head.Seq, cs.ERR_CODE_FAILED)
		return
	}

	var res cs.CsCmdSceneEnterRes
	res.MapId = map_id
	res.SceneId = scene_id

	player.SendMsg(cs.CS_CMD_SCENE_ENTER_RES, head.Seq, &res)
}

func HandleCsCmdSceneLeaveReq(player *model.Player, head *cs.CsHead, body *proto.Message) {
	if player == nil || head == nil || body == nil {
		return
	}

	req := (*body).(*cs.CsCmdSceneLeaveReq)
	if req == nil {
		fmt.Println("HandleCsCmdSceneLeaveReq req is nil")
		player.SendErrCode(cs.CS_CMD_SCENE_LEAVE_RES, head.Seq, cs.ERR_CODE_FAILED)
		return
	}

	map_id := req.GetMapId()
	scene_id := req.GetSceneId()

	if !scene.SceneMgrIn.LeavePlayer(player, map_id, scene_id) {
		player.SendErrCode(cs.CS_CMD_SCENE_LEAVE_RES, head.Seq, cs.ERR_CODE_FAILED)
		return
	}

	var res cs.CsCmdSceneLeaveRes
	res.MapId = map_id
	res.SceneId = scene_id

	player.SendMsg(cs.CS_CMD_SCENE_LEAVE_RES, head.Seq, &res)
}

func HandleCsCmdSceneMoveReq(player *model.Player, head *cs.CsHead, body *proto.Message) {
	if player == nil || head == nil || body == nil {
		return
	}

	req := (*body).(*cs.CsCmdSceneMoveReq)
	if req == nil {
		fmt.Println("HandleCsCmdSceneMoveReq req is nil")
		player.SendErrCode(cs.CS_CMD_SCENE_MOVE_RES, head.Seq, cs.ERR_CODE_FAILED)
		return
	}

	map_id := req.GetMapId()
	scene_id := req.GetSceneId()
	x := req.GetPosX()
	y := req.GetPosY()

	scene := scene.SceneMgrIn.GetScene(map_id, scene_id)
	if scene == nil {
		player.SendErrCode(cs.CS_CMD_SCENE_MOVE_RES, head.Seq, cs.ERR_CODE_FAILED)
		return
	}

	if !scene.ChangeGrid(player.GetAccountID(), float32(x)/100, float32(y)/100) {
		player.SendErrCode(cs.CS_CMD_SCENE_MOVE_RES, head.Seq, cs.ERR_CODE_FAILED)
		return
	}

	var res cs.CsCmdSceneMoveRes
	res.MapId = map_id
	res.SceneId = scene_id

	player.SendMsg(cs.CS_CMD_SCENE_MOVE_RES, head.Seq, &res)
}
