package handle

import (
	"client/agent"
	"fmt"
	"jserver/src/common/message"
)

type HandleFunc func(*agent.AgentClient, *message.Message)

var HandleMap map[uint32]HandleFunc

func init() {
	fmt.Println("HandleMap init")
	HandleMap = make(map[uint32]HandleFunc)
}

func RegisterHandle(id uint32, handle HandleFunc) {
	if _, ok := HandleMap[id]; ok {
		panic("handle already registered")
	}
	HandleMap[id] = handle
}

func Handle(p *agent.AgentClient, msg *message.Message) {
	if handle, ok := HandleMap[msg.Head.Id]; ok {
		handle(p, msg)
	}
}

func HandleReceiveMessage(c *agent.AgentClient) {
	msg, err := c.ReceiveMessage()
	if err != nil {
		fmt.Println("Failed to receive message: ", err)
		return
	}

	Handle(c, msg)
}
