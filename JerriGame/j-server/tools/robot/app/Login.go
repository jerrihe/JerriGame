package app

import (
	"fmt"
	handle "robot/handle"
	model "robot/model"

	"jserver/src/protocol/cs"

	"google.golang.org/protobuf/proto"
)

func init() {
	fmt.Println("app init")

	handle.RegisterHandleFunc(cs.CS_CMD_LOGIN_RES, HandleCsCmdLoginRes)
	handle.RegisterHandleFunc(cs.CS_CMD_CREATE_ACCOUNT_REQ, HandleCsCmdCreateRoleRes)
	handle.RegisterHandleFunc(cs.CS_CMD_NTF_KICK_ACCOUNT, HandleCsCmdNtfKickAccount)
	handle.RegisterHandleFunc(cs.CS_CMD_NTF_ERROR_CODE, HandleCsCmdNtfErrorCode)
	handle.RegisterHandleFunc(cs.CS_CMD_LOGIN_OUT_RES, HandleCsCmdLoginOutRes)
}

func HandleCsCmdLoginRes(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	var res = (*msg).(*cs.CsCmdLoginRes)

	agent.SetAccountID(res.AccountId)
	agent.SetState(model.ActiveEnumLoginSuccess)

	fmt.Println("LoginRes:", res)
}

func HandleCsCmdCreateRoleRes(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	fmt.Println("app HandleCsCmdCreateRoleRes")

	var res = (*msg).(*cs.CsCmdCreateAccountRes)
	// 已登录状态
	agent.SetAccountID(res.AccountId)
	agent.SetState(model.ActiveEnumCreateRoleSuccess)

	fmt.Println("CreateRoleRes:", res)
}

func HandleCsCmdNtfKickAccount(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	fmt.Println("app HandleKick")

	var res = (*msg).(*cs.CsCmdNtfKickAccount)

	fmt.Println("HandleKick:", res)

	// 离线状态
	agent.SetState(model.ActiveEnumKick)
	// 需要关闭连接
}

func HandleCsCmdNtfErrorCode(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	fmt.Println("app HandleErrorCode")

	var ntf = (*msg).(*cs.CsCmdNtfErrorCode)

	fmt.Println("HandleErrorCode:", ntf.Cmd, ntf.ErrCode)

	handle.HandleErrCode(agent, cs.CS_CMD(ntf.Cmd), cs.ERR_CODE(ntf.ErrCode))

}

func HandleCsCmdLoginOutRes(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	fmt.Println("app HandleLoginOut")

	var res = (*msg).(*cs.CsCmdLoginOutRes)

	fmt.Println("HandleLoginOut:", res)

	// 登出成功
	agent.SetState(model.ActiveEnumLoginOutSuccess)
}
