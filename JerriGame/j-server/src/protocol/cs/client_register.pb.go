// Code generated by register. DO NOT EDIT.

package cs

import "google.golang.org/protobuf/reflect/protoreflect"

import "google.golang.org/protobuf/proto"

type EnmCmdValue CS_CMD

type IMessage interface {
	Reset()
	String() string
	ProtoMessage()
	ProtoReflect() protoreflect.Message
	Descriptor() ([]byte, []int)
}

type newMessage func () proto.Message

var imessageMap map[int32] newMessage
var cmdReqResMap map[int32] int32

func NewClientMessage(CmdId int32) proto.Message {
	v, ok := imessageMap[CmdId]
	if ! ok {
		return nil
	}
	return v()
}

func GetResCmdId(reqCmdId int32) int32{
	v, ok := cmdReqResMap[reqCmdId]
	if ! ok {
		return 0
	}
	return v
}

func newCsCmdSceneNtfMove() proto.Message {
	return &CsCmdSceneNtfMove{}
}

func newCsCmdSceneMoveRes() proto.Message {
	return &CsCmdSceneMoveRes{}
}

func newCsCmdLoginReq() proto.Message {
	return &CsCmdLoginReq{}
}

func newCsCmdCreateAccountRes() proto.Message {
	return &CsCmdCreateAccountRes{}
}

func newCsCmdSceneNtfLeave() proto.Message {
	return &CsCmdSceneNtfLeave{}
}

func newCsCmdSceneNtfEnter() proto.Message {
	return &CsCmdSceneNtfEnter{}
}

func newCsCmdLoginOutRes() proto.Message {
	return &CsCmdLoginOutRes{}
}

func newCsCmdCreateAccountReq() proto.Message {
	return &CsCmdCreateAccountReq{}
}

func newCsCmdSceneLeaveReq() proto.Message {
	return &CsCmdSceneLeaveReq{}
}

func newCsCmdSceneLeaveRes() proto.Message {
	return &CsCmdSceneLeaveRes{}
}

func newCsCmdLoginRes() proto.Message {
	return &CsCmdLoginRes{}
}

func newCsCmdLoginOutReq() proto.Message {
	return &CsCmdLoginOutReq{}
}

func newCsCmdSceneEnterRes() proto.Message {
	return &CsCmdSceneEnterRes{}
}

func newCsCmdNtfKickAccount() proto.Message {
	return &CsCmdNtfKickAccount{}
}

func newCsCmdNtfErrorCode() proto.Message {
	return &CsCmdNtfErrorCode{}
}

func newCsCmdSceneEnterReq() proto.Message {
	return &CsCmdSceneEnterReq{}
}

func newCsCmdSceneMoveReq() proto.Message {
	return &CsCmdSceneMoveReq{}
}

func init() {
	imessageMap = make(map[int32] newMessage)
	imessageMap[1] = newCsCmdLoginReq
	imessageMap[6] = newCsCmdCreateAccountRes
	imessageMap[206] = newCsCmdSceneNtfLeave
	imessageMap[207] = newCsCmdSceneNtfMove
	imessageMap[209] = newCsCmdSceneMoveRes
	imessageMap[4] = newCsCmdLoginOutRes
	imessageMap[5] = newCsCmdCreateAccountReq
	imessageMap[203] = newCsCmdSceneLeaveReq
	imessageMap[204] = newCsCmdSceneLeaveRes
	imessageMap[205] = newCsCmdSceneNtfEnter
	imessageMap[2] = newCsCmdLoginRes
	imessageMap[3] = newCsCmdLoginOutReq
	imessageMap[202] = newCsCmdSceneEnterRes
	imessageMap[7] = newCsCmdNtfKickAccount
	imessageMap[101] = newCsCmdNtfErrorCode
	imessageMap[201] = newCsCmdSceneEnterReq
	imessageMap[208] = newCsCmdSceneMoveReq

	cmdReqResMap = make(map[int32]int32)
	cmdReqResMap[201] = 202
	cmdReqResMap[208] = 209
	cmdReqResMap[1] = 2
	cmdReqResMap[5] = 6
	cmdReqResMap[203] = 204
	cmdReqResMap[3] = 4
}

func GetCmdValueByMsg(msg interface{}) EnmCmdValue {
	if nil == msg {
		return EnmCmdValue(0)
	}
	switch msg.(type) {
		case CsCmdLoginRes:
			return EnmCmdValue(2)
		case *CsCmdLoginRes:
			return EnmCmdValue(2)
		case CsCmdLoginOutReq:
			return EnmCmdValue(3)
		case *CsCmdLoginOutReq:
			return EnmCmdValue(3)
		case CsCmdSceneEnterRes:
			return EnmCmdValue(202)
		case *CsCmdSceneEnterRes:
			return EnmCmdValue(202)
		case CsCmdNtfKickAccount:
			return EnmCmdValue(7)
		case *CsCmdNtfKickAccount:
			return EnmCmdValue(7)
		case CsCmdNtfErrorCode:
			return EnmCmdValue(101)
		case *CsCmdNtfErrorCode:
			return EnmCmdValue(101)
		case CsCmdSceneEnterReq:
			return EnmCmdValue(201)
		case *CsCmdSceneEnterReq:
			return EnmCmdValue(201)
		case CsCmdSceneMoveReq:
			return EnmCmdValue(208)
		case *CsCmdSceneMoveReq:
			return EnmCmdValue(208)
		case CsCmdLoginReq:
			return EnmCmdValue(1)
		case *CsCmdLoginReq:
			return EnmCmdValue(1)
		case CsCmdCreateAccountRes:
			return EnmCmdValue(6)
		case *CsCmdCreateAccountRes:
			return EnmCmdValue(6)
		case CsCmdSceneNtfLeave:
			return EnmCmdValue(206)
		case *CsCmdSceneNtfLeave:
			return EnmCmdValue(206)
		case CsCmdSceneNtfMove:
			return EnmCmdValue(207)
		case *CsCmdSceneNtfMove:
			return EnmCmdValue(207)
		case CsCmdSceneMoveRes:
			return EnmCmdValue(209)
		case *CsCmdSceneMoveRes:
			return EnmCmdValue(209)
		case CsCmdLoginOutRes:
			return EnmCmdValue(4)
		case *CsCmdLoginOutRes:
			return EnmCmdValue(4)
		case CsCmdCreateAccountReq:
			return EnmCmdValue(5)
		case *CsCmdCreateAccountReq:
			return EnmCmdValue(5)
		case CsCmdSceneLeaveReq:
			return EnmCmdValue(203)
		case *CsCmdSceneLeaveReq:
			return EnmCmdValue(203)
		case CsCmdSceneLeaveRes:
			return EnmCmdValue(204)
		case *CsCmdSceneLeaveRes:
			return EnmCmdValue(204)
		case CsCmdSceneNtfEnter:
			return EnmCmdValue(205)
		case *CsCmdSceneNtfEnter:
			return EnmCmdValue(205)
	}
	return EnmCmdValue(0)
}

