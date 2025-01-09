package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	gnet "github.com/walkon/wsgnet"
	"google.golang.org/protobuf/proto"

	"jserver/src/protocol/cs"
	"jserver/src/protocol/ss"
)

func SendToRouter(conn gnet.Conn, cmd ss.SS_CMD, server_id int32, target_server_id int32, body interface{}) {
	if body == nil || conn == nil {
		return
	}

	msg := body.(proto.Message)
	bodyData, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("Error serializing body struct:", err)
		return
	}
	// 创建消息
	var head = &ss.SsHead{Cmd: int32(cmd), ServerId: server_id, TargetServerId: target_server_id}
	// 消息包
	var ssPkg = &ss.SsPkg{Head: head, Body: bodyData}
	// 最终的协议包
	protocolPacket, err := SerializeMessage(ssPkg)
	if err != nil {
		fmt.Println("Error serializing message:", err, " cmd:", cmd)
		return
	}

	// 发送消息
	conn.AsyncWrite(protocolPacket, nil)

	fmt.Println("Sent message to server: ", target_server_id, " msg: ", body, "len", len(protocolPacket))
}

func SendToClientMsg(conn gnet.Conn, cmd int32, seqId int32, body interface{}) {
	if conn == nil || body == nil {
		return
	}

	msg := body.(proto.Message)

	var head = &cs.CsHead{Cmd: cmd, Seq: seqId, Ret: 0}
	bodyData, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("Error serializing body struct:", err)
		return
	}
	// 消息包
	var csPkg = &cs.CsPkg{Head: head, Body: bodyData}
	// 最终的协议包
	protocolPacket, err := SerializeMessage(csPkg)
	if err != nil {
		fmt.Println("Error serializing message:", err, " cmd:", cmd)
		return
	}

	// 发送消息
	conn.AsyncWrite(protocolPacket, nil)
}

func SendErrCodeToClient(conn gnet.Conn, cmd cs.CS_CMD, seqId int32, errCode cs.ERR_CODE) {
	if conn == nil {
		return
	}

	var head = &cs.CsHead{Cmd: int32(cs.CS_CMD_NTF_ERROR_CODE), Seq: seqId, Ret: int32(errCode)}
	var body = &cs.CsCmdNtfErrorCode{}
	body.Cmd = int32(cmd)
	body.ErrCode = int32(errCode)

	bodyData, err := proto.Marshal(body)
	if err != nil {
		fmt.Println("Error serializing body struct:", err)
		return
	}

	// 消息包
	var csPkg = &cs.CsPkg{Head: head, Body: bodyData}
	// 最终的协议包
	protocolPacket, err := SerializeMessage(csPkg)
	if err != nil {
		fmt.Println("Error serializing message:", err, " cmd:", cmd)
		return
	}

	// 发送消息
	conn.AsyncWrite(protocolPacket, nil)
}

func SerializeMessage(pkg proto.Message) ([]byte, error) {
	// 序列化消息头
	pKgData, err := proto.Marshal(pkg)
	if err != nil {
		fmt.Println("Error serializing head struct:", err)
		return nil, err
	}

	// 包体总长度 = 4字节包体长度 + 包体信息
	totalLen := uint32(4 + len(pKgData)) // 4 字节的包体总长度 + 包信息

	// 拼接协议包
	var buffer bytes.Buffer
	// 写入包体总长度（4 字节）
	err = binary.Write(&buffer, binary.BigEndian, totalLen)
	if err != nil {
		fmt.Println("Failed to write total length:", err)
		return nil, err
	}

	// 写入包体数据
	_, err = buffer.Write(pKgData)
	if err != nil {
		log.Fatalf("Failed to write body data: %v", err)
	}

	// 最终的协议包
	protocolPacket := buffer.Bytes()

	return protocolPacket, nil
}

func CSParseProtocolPacket(packet []byte) (uint32, *cs.CsPkg, *proto.Message, error) {
	// 读取包体总长度（4 字节）
	totalLen := binary.BigEndian.Uint32(packet[:4])
	// fmt.Printf("Total Length: %d\n", totalLen)

	if totalLen > uint32(len(packet)) {
		return totalLen, nil, nil, errors.New("error: invalid packet length")
	}

	// 读取包体数据
	body := packet[4:]

	var csPkg cs.CsPkg
	err := proto.Unmarshal(body, &csPkg)
	if err != nil {
		return totalLen, nil, nil, err
	}

	cmd := csPkg.Head.Cmd
	bodyMsg := cs.NewClientMessage(cmd)
	err = proto.Unmarshal(csPkg.Body, bodyMsg)
	if err != nil {
		fmt.Println("Error deserializing body struct:", err)
		return totalLen, nil, nil, err
	}
	return totalLen, &csPkg, &bodyMsg, nil
}

func CSParseProtocolPacket1(packet []byte) (uint32, *cs.CsPkg, error) {
	// 读取包体总长度（4 字节）
	totalLen := binary.BigEndian.Uint32(packet[:4])
	fmt.Printf("Total Length: %d\n", totalLen)

	if totalLen > uint32(len(packet)) {
		return totalLen, nil, errors.New("error: invalid packet length")
	}

	// 读取包体数据
	body := packet[4:]

	var csPkg cs.CsPkg
	err := proto.Unmarshal(body, &csPkg)
	if err != nil {
		return totalLen, nil, err
	}
	fmt.Printf("Body: %s\n", csPkg.String())
	return totalLen, &csPkg, nil
}

// 读取包体总长度（4 字节）
func SSParseProtocolPacket(packet []byte) (uint32, *ss.SsPkg, proto.Message, error) {
	totalLen := binary.BigEndian.Uint32(packet[:4])
	// fmt.Printf("Total Length: %d\n", totalLen)

	if totalLen > uint32(len(packet)) {
		return totalLen, nil, nil, errors.New("error: invalid packet length")
	}

	// 读取包体数据
	body := packet[4:]
	var ssPkg ss.SsPkg
	err := proto.Unmarshal(body, &ssPkg)
	if err != nil {
		fmt.Println("Error deserializing csPkg struct:", err)
		return totalLen, nil, nil, err
	}

	cmd := ssPkg.Head.Cmd
	bodyMsg := ss.NewServiceMessage(cmd)
	err = proto.Unmarshal(ssPkg.Body, bodyMsg)
	if err != nil {
		fmt.Println("Error deserializing body struct:", err)
		return totalLen, nil, nil, err
	}
	return totalLen, &ssPkg, bodyMsg, nil
}
