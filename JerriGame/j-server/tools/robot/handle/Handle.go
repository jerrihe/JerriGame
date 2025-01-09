package handle

import (
	"fmt"
	"jserver/src/protocol/cs"
	model "robot/model"

	"google.golang.org/protobuf/proto"
)

type HandleFunc func(*model.Agent, *cs.CsHead, *proto.Message)
type HandleErrCodeFunc func(*model.Agent, cs.CS_CMD, cs.ERR_CODE)

var HandleFuncMap map[cs.CS_CMD]HandleFunc
var HandleErrCodeMap map[cs.CS_CMD]HandleErrCodeFunc

func init() {
	fmt.Println("Handle init")
	HandleFuncMap = make(map[cs.CS_CMD]HandleFunc)
	HandleErrCodeMap = make(map[cs.CS_CMD]HandleErrCodeFunc)
}

func RegisterHandleFunc(msg_id cs.CS_CMD, handle_func HandleFunc) {
	if _, ok := HandleFuncMap[msg_id]; ok {
		panic("msg_id already registered")

	} else {
		HandleFuncMap[msg_id] = handle_func
	}
}

func GetHandleFunc(msg_id cs.CS_CMD) HandleFunc {
	return HandleFuncMap[msg_id]
}

func HandleMsg(agent *model.Agent, head *cs.CsHead, msg *proto.Message) {
	handle_func := GetHandleFunc(cs.CS_CMD(head.Cmd))
	if handle_func != nil {
		handle_func(agent, head, msg)
	}
}

func RegisterHandleErrCode(msg_id cs.CS_CMD, handle_func HandleErrCodeFunc) {
	if _, ok := HandleErrCodeMap[msg_id]; ok {
		panic("msg_id already registered")

	} else {
		HandleErrCodeMap[msg_id] = handle_func
	}
}

func GetHandleErrCode(msg_id cs.CS_CMD) HandleErrCodeFunc {
	return HandleErrCodeMap[msg_id]
}

func HandleErrCode(agent *model.Agent, cmd cs.CS_CMD, errCode cs.ERR_CODE) {
	handle_func := GetHandleErrCode(cmd)
	if handle_func != nil {
		handle_func(agent, cmd, errCode)
	}
}
