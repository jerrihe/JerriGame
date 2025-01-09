package app

import (
	"fmt"
	"jserver/src/protocol/cs"
	handle "robot/handle"
	model "robot/model"
)

func init() {
	handle.RegisterHandleErrCode(cs.CS_CMD_LOGIN_RES, HandleLoginResErrCode)
	handle.RegisterHandleErrCode(cs.CS_CMD_CREATE_ACCOUNT_REQ, HandleCreateRoleResErrCode)
}

func HandleLoginResErrCode(agent *model.Agent, cmd cs.CS_CMD, errCode cs.ERR_CODE) {
	fmt.Println("app HandleLoginResErrCode")

	fmt.Println("HandleLoginResErrCode:", cmd, errCode)
	switch errCode {
	case cs.ERR_CODE_LOGIN_NEED_CREATE:
		// 需要创建角色
		fmt.Println("need create role")
		agent.SetState(model.ActiveEnumCreateRole)
	}
}

func HandleCreateRoleResErrCode(agent *model.Agent, cmd cs.CS_CMD, errCode cs.ERR_CODE) {
	fmt.Println("app HandleCreateRoleResErrCode")

	fmt.Println("HandleCreateRoleResErrCode:", cmd, errCode)
	switch errCode {
	case cs.ERR_CODE_CREATE_ACCOUNT_FAILED:
		// 创建角色失败
		fmt.Println("create role failed")
		agent.SetState(model.ActiveEnumKick)
	case cs.ERR_CODE_CREATE_EXITE_ACCOUNT:
		// 角色已存在
		fmt.Println("role already exist")
		agent.SetState(model.ActiveEnumLogin)
	}
}
