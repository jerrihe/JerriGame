package agent

import (
	"encoding/json"
	"fmt"
	"jserver/src/common/message"

	"net"
	"sync"
)

type AgentClient struct {
	serverAddress string     // 服务端地址
	conn          net.Conn   // 客户端连接
	lock          sync.Mutex // 保证线程安全

	UserID   string
	PlatType string

	account_id uint64
	role_id    uint64

	is_login bool
}

func (c *AgentClient) SetAccountId(account_id uint64) {
	c.account_id = account_id
}

func (c *AgentClient) GetAccountId() uint64 {
	return c.account_id
}

func (c *AgentClient) SetRoleId(role_id uint64) {
	c.role_id = role_id
}

func (c *AgentClient) GetRoleId() uint64 {

	return c.role_id
}

func (c *AgentClient) SetIsLogin(is_login bool) {
	c.is_login = is_login
}

func (c *AgentClient) GetIsLogin() bool {
	return c.is_login
}

func (c *AgentClient) Connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var err error
	c.conn, err = net.Dial("tcp", c.serverAddress)
	if err != nil {
		return err
	}

	fmt.Printf("Connected to server at %s\n", c.serverAddress)
	return nil
}

func (c *AgentClient) SendMessage(cmd uint32, msgContent interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conn == nil {
		return net.ErrClosed
	}

	bytes, err := json.Marshal(msgContent)
	if err != nil {
		fmt.Println("Error serializing struct:", err)
		return nil
	}

	sendmsg := message.NewMessage(cmd, bytes)

	buf, err := sendmsg.Encode()
	if err != nil {
		return nil
	}
	_, err = c.conn.Write(buf)
	if err != nil {
		return err
	}

	fmt.Println("Sent message to server: ", msgContent)

	return nil
}

func (c *AgentClient) ReceiveMessage() (*message.Message, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	headbuf := make([]byte, 8)
	_, err := c.conn.Read(headbuf)
	if err != nil {
		return nil, err
	}

	head, err := message.DecodeHead(headbuf)
	if err != nil {
		return nil, err
	}

	content_buf := make([]byte, head.ContentLength)
	n, err := c.conn.Read(content_buf)
	if err != nil {
		return nil, err
	}

	if uint32(n) != head.ContentLength {
		return nil, fmt.Errorf("message length mismatch")
	}

	return message.NewMessage(head.Id, content_buf), nil
}

func (c *AgentClient) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conn == nil {
		return net.ErrClosed
	}

	err := c.conn.Close()
	if err != nil {
		return err
	}

	c.conn = nil
	return nil
}

func NewAgentClient(serverAddress string, user string, plat string) *AgentClient {
	return &AgentClient{serverAddress: serverAddress,
		UserID:   user,
		PlatType: plat}
}
