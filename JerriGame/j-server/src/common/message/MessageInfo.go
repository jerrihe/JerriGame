package message

import (
	"encoding/json"
	"fmt"
)

type HeadInfo struct {
	Id             uint32
	ServerID       uint32
	TargetServerID uint32
	ContentLength  uint32
}

type ServerMessage struct {
	HeadInfo
	Content []byte
}

func NewServerMessage(id uint32, serverID uint32, targetServerID uint32, content interface{}) *ServerMessage {
	contentbUF, err := json.Marshal(content)
	if err != nil {
		return nil
	}
	return &ServerMessage{
		HeadInfo: HeadInfo{
			Id:             id,
			ServerID:       serverID,
			TargetServerID: targetServerID,
			ContentLength:  uint32(len(contentbUF)),
		},
		Content: contentbUF,
	}
}

func (m *ServerMessage) Encode() ([]byte, error) {
	if m.ContentLength != uint32(len(m.Content)) {
		return nil, fmt.Errorf("message length mismatch")
	}

	if m.Id == 0 {
		return nil, fmt.Errorf("message id is 0")
	}

	// 16 字节的消息头部
	// 4 字节的消息 ID
	tolength := 16 + m.ContentLength
	buf := make([]byte, tolength)
	buf[0] = byte(m.Id)
	buf[1] = byte(m.Id >> 8)
	buf[2] = byte(m.Id >> 16)
	buf[3] = byte(m.Id >> 24)

	// 4 字节的消息 ServerID
	buf[4] = byte(m.ServerID)
	buf[5] = byte(m.ServerID >> 8)
	buf[6] = byte(m.ServerID >> 16)
	buf[7] = byte(m.ServerID >> 24)

	// 4 字节的消息 TargetServerID
	buf[8] = byte(m.TargetServerID)
	buf[9] = byte(m.TargetServerID >> 8)
	buf[10] = byte(m.TargetServerID >> 16)
	buf[11] = byte(m.TargetServerID >> 24)

	//  4 字节的消息长度
	buf[12] = byte(m.ContentLength)
	buf[13] = byte(m.ContentLength >> 8)
	buf[14] = byte(m.ContentLength >> 16)
	buf[15] = byte(m.ContentLength >> 24)

	copy(buf[16:], m.Content)
	return buf, nil
}

func DecodeServerMessage(buf []byte) (*ServerMessage, uint32, error) {
	if len(buf) < 16 {
		return nil, 0, fmt.Errorf("message length is too short")
	}

	id := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
	serverID := uint32(buf[4]) | uint32(buf[5])<<8 | uint32(buf[6])<<16 | uint32(buf[7])<<24
	targetServerID := uint32(buf[8]) | uint32(buf[9])<<8 | uint32(buf[10])<<16 | uint32(buf[11])<<24
	length := uint32(buf[12]) | uint32(buf[13])<<8 | uint32(buf[14])<<16 | uint32(buf[15])<<24

	if id == 0 {
		return nil, 0, fmt.Errorf("message id is 0")
	}

	if uint32(len(buf)) < length+16 {
		fmt.Println("message length mismatch", length, len(buf))
		return nil, 0, fmt.Errorf("message length mismatch")
	}

	content := buf[16:]
	return &ServerMessage{
		HeadInfo: HeadInfo{
			Id:             id,
			ServerID:       serverID,
			TargetServerID: targetServerID,
			ContentLength:  length,
		},
		Content: content,
	}, length + 16, nil
}
