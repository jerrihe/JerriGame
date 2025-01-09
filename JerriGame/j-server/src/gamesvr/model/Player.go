package model

import (
	"jserver/src/common/message"
	"jserver/src/protocol/cs"

	gnet "github.com/walkon/wsgnet"
)

type PlayerPos struct {
	X float32
	Y float32
}

type PlayerSceneInfo struct {
	// 场景ID
	SceneID uint32
	// mapID
	MapID uint32

	Pos    PlayerPos
	OldPos PlayerPos

	GridID  int32
	OldGrid int32

	// 看见的玩家
	SeePlayers map[uint64]bool
}

type Player struct {
	User      string
	Platform  string
	AccountID uint64

	Conn    gnet.Conn
	ConnIdx int64

	is_login bool

	SceneInfo PlayerSceneInfo
}

func (p *Player) Init(conn gnet.Conn, conn_idx int64) {
	p.Conn = conn
	p.ConnIdx = conn_idx
}

func (p *Player) InitLogin(user string, platform string) {
	p.User = user
	p.Platform = platform
	// p.AccountID = account_id
	p.is_login = false
}

// 场景进入初始化
func (p *Player) InitSceneInfo(scene_id uint32, map_id uint32, x float32, y float32, grid_id int32) {
	p.SceneInfo.SceneID = scene_id
	p.SceneInfo.MapID = map_id
	p.SceneInfo.Pos.X = x
	p.SceneInfo.Pos.Y = y
	p.SceneInfo.GridID = grid_id
	p.SceneInfo.SeePlayers = make(map[uint64]bool)
	p.SceneInfo.OldPos.X = x
	p.SceneInfo.OldPos.Y = y
	p.SceneInfo.OldGrid = grid_id
}

// 场景退出
func (p *Player) ClearSceneInfo() {
	p.SceneInfo.SceneID = 0
	p.SceneInfo.MapID = 0
	p.SceneInfo.Pos.X = -1
	p.SceneInfo.Pos.Y = -1
	p.SceneInfo.GridID = -1
	p.SceneInfo.SeePlayers = make(map[uint64]bool)
	p.SceneInfo.OldPos.X = -1
	p.SceneInfo.OldPos.Y = -1
	p.SceneInfo.OldGrid = -1
}

// 有玩家进入视野
func (p *Player) AddSeePlayer(account_id uint64) {
	p.SceneInfo.SeePlayers[account_id] = true

	// 发送进入视野消息
}

// 有玩家离开视野
func (p *Player) DelSeePlayer(account_id uint64) {
	if _, ok := p.SceneInfo.SeePlayers[account_id]; !ok {
		return
	}
	delete(p.SceneInfo.SeePlayers, account_id)

	// 发送离开视野消息
}

func (p *Player) GetConn() gnet.Conn {
	return p.Conn
}

func (p *Player) GetUser() string {
	return p.User
}

func (p *Player) GetPlatform() string {
	return p.Platform
}

func (p *Player) GetAccountID() uint64 {
	return p.AccountID
}

func (p *Player) SetAccountID(account_id uint64) {
	p.AccountID = account_id
}

func (p *Player) IsLogin() bool {
	return p.is_login
}

func (p *Player) GetConnIdx() int64 {
	return p.ConnIdx
}

func (p *Player) SetConnIdx(conn_idx int64) {
	p.ConnIdx = conn_idx
}

func (p *Player) SetConn(conn gnet.Conn) {
	p.Conn = conn
}

func (p *Player) GetUserAndPlatform() string {
	return p.User + "_" + p.Platform
}

func (p *Player) SetIsLogin(is_login bool) {
	p.is_login = is_login
}

func (p *Player) SendMsg(cmd cs.CS_CMD, seqId int32, body interface{}) {
	message.SendToClientMsg(p.Conn, int32(cmd), seqId, body)
}

func (p *Player) SendErrCode(cmd cs.CS_CMD, seqId int32, errCode cs.ERR_CODE) {
	message.SendErrCodeToClient(p.Conn, cmd, seqId, errCode)
}
