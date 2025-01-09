package model

import (
	"time"

	gnet "github.com/walkon/wsgnet"
)

type Agent struct {
	User      string
	Platform  string
	AccountID uint64
	Conn      gnet.Conn

	State        int
	ConnectState int
	ConnIdx      uint64

	MapId   uint32
	SceneId uint32

	PosX float32
	PosY float32

	// 当前执行的活动类型
	ActiveType uint32
	// 当前执行节点ID
	ActiveNodeID uint32
	// 当前执行子节点ID
	ActiveSubNodeIndex uint32
}

func (a *Agent) Init(user string, plat_type string, connIdx uint64, conn gnet.Conn) {
	a.User = user
	a.Platform = plat_type
	a.State = ActiveEnumLogin
	a.ConnectState = ConnectNone
	a.ConnIdx = connIdx
	a.Conn = conn

	a.ActiveType = ACTIVE_TYPE_SCENE
	a.ActiveNodeID = 1
	a.ActiveSubNodeIndex = 0
}

func (a *Agent) SetSceneInfo(map_id uint32, scene_id uint32, x float32, y float32) {
	a.MapId = map_id
	a.SceneId = scene_id
	a.PosX = x
	a.PosY = y
}

func (a *Agent) ClearSceneInfo() {
	a.MapId = 0
	a.SceneId = 0
	a.PosX = 0
	a.PosY = 0
}

func (a *Agent) SetPos(x float32, y float32) {
	a.PosX = x
	a.PosY = y
}

func (a *Agent) GetPos() (float32, float32) {
	return a.PosX, a.PosY
}

func (a *Agent) Tick1s(now time.Time) {

}

func (a *Agent) SetConnectState(state int) {
	a.ConnectState = state
}

func (a *Agent) GetConnectState() int {
	return a.ConnectState
}

func (a *Agent) SetConn(conn gnet.Conn) {
	a.Conn = conn
}

func (a *Agent) SetAccountID(account_id uint64) {
	a.AccountID = account_id
}

func (a *Agent) GetAccountID() uint64 {
	return a.AccountID
}

func (a *Agent) GetUser() string {
	return a.User
}

func (a *Agent) GetPlatform() string {
	return a.Platform
}

func (a *Agent) GetConn() gnet.Conn {
	return a.Conn
}

func (a *Agent) SetState(state int) {
	a.State = state
}

func (a *Agent) GetState() int {
	return a.State
}
