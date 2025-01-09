package app

import (
	agent "client/agent"
	pool "client/workpool"
	"encoding/json"
	"fmt"
	message "jserver/src/common/message"
	"time"

	handle "client/handle"
)

var ServeLoginNode map[int][]int

func InitApp() {
	fmt.Println("InitApp")

	handle.RegisterHandle(message.CS_RES_LOGIN, HandleLoginRes)
	handle.RegisterHandle(message.CS_RES_CREATE_ROLE, HandleCreateRoleRes)
}

func RunAgentClient(serverAddress string, user string, plat string) {
	client := agent.NewAgentClient(serverAddress, user, plat)
	err := client.Connect()
	if err != nil {
		fmt.Println("Failed to connect to server: ", err)
		return
	}

	defer client.Close()

	for {
		if client.GetIsLogin() {
			fmt.Println("Client is login")
		} else {
			// 登录
			HandleLoginReq(client)
		}
		time.Sleep(1 * time.Second) // Add a sleep interval to prevent 100% CPU usage
	}
}

func HandleLoginReq(c *agent.AgentClient) {
	cmd := uint32(message.CS_REQ_LOGIN)
	var req message.CSLoginReq
	req.User = c.UserID
	req.PlatType = c.PlatType

	//send := func() {
	c.SendMessage(cmd, req)
	handle.HandleReceiveMessage(c)
	//}

	// pool.Worker.AddTask(pool.Task{ID: int(cmd), Task: send})
	fmt.Println("HandleLoginReq")
}

func HandleLoginRes(c *agent.AgentClient, msg *message.Message) {
	fmt.Println("HandleLoginRes")
	var res message.CSLoginRes
	err := json.Unmarshal(msg.Content, &res)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
	}

	if res.Result == 0 {
		c.SetIsLogin(true)
		c.SetAccountId(res.AccID)
	} else {
		fmt.Println("Login failed")
		// 创建角色
		send := func() {
			HandleCreateRoleReq(c)
		}
		pool.Worker.AddTask(pool.Task{ID: int(message.CS_REQ_CREATE_ROLE), Task: send})
	}
	fmt.Println("CSLoginRes:", res)
}

func HandleCreateRoleReq(c *agent.AgentClient) {
	cmd := uint32(message.CS_REQ_CREATE_ROLE)
	var req message.CSCreateRoleReq
	req.RoleName = "test"
	req.User = c.UserID
	req.PlatType = c.PlatType

	// 询问平台创号

	//send := func() {
	c.SendMessage(cmd, req)
	handle.HandleReceiveMessage(c)
	//}

	//pool.Worker.AddTask(pool.Task{ID: int(cmd), Task: send})
	fmt.Println("HandleCreateRoleReq")
}

func HandleCreateRoleRes(c *agent.AgentClient, msg *message.Message) {
	fmt.Println("HandleCreateRoleRes")
	var res message.CSCreateRoleRes
	err := json.Unmarshal(msg.Content, &res)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
	}

	if res.Result == 0 {
		c.SetRoleId(res.AccID)
		fmt.Println("Create role success")
		c.SetIsLogin(true)
	} else {
		fmt.Println("Create role failed")
	}
}
