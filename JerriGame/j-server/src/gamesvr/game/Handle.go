package game

import (
	"fmt"
	"jserver/src/protocol/cs"
	"jserver/src/protocol/ss"

	model "gamesvr/model"

	proto "google.golang.org/protobuf/proto"
)

type HandleFunc func(*model.Player, *cs.CsHead, *proto.Message)
type HandleSSFunc func(*ss.SsHead, *proto.Message)
type HandleErrCode func(cmd cs.CS_CMD, errCode int32)

var HandleMap map[cs.CS_CMD]HandleFunc
var HandleSSMap map[ss.SS_CMD]HandleSSFunc
var HandleErrCodeMap map[cs.CS_CMD]HandleErrCode

func init() {
	fmt.Println("HandleMap init")
	HandleMap = make(map[cs.CS_CMD]HandleFunc)
	HandleSSMap = make(map[ss.SS_CMD]HandleSSFunc)
}

func RegisterHandle(id cs.CS_CMD, handle HandleFunc) {
	if _, ok := HandleMap[id]; ok {
		panic("handle already registered")
	}
	HandleMap[id] = handle
}

func HandleCs(p *model.Player, head *cs.CsHead, body *proto.Message) {
	if handle, ok := HandleMap[cs.CS_CMD(head.GetCmd())]; ok {
		handle(p, head, body)
	}
}

func RegisterHandleSS(id ss.SS_CMD, handle HandleSSFunc) {
	if _, ok := HandleSSMap[id]; ok {
		panic("handle already registered")
	}
	HandleSSMap[id] = handle
}

func HandleSS(head *ss.SsHead, body *proto.Message) {
	if handle, ok := HandleSSMap[ss.SS_CMD(head.GetCmd())]; ok {
		handle(head, body)
	}
}
