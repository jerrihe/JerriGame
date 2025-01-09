package model

import (
	"fmt"
	"sync"
	"time"

	gnet "github.com/walkon/wsgnet"
)

type playerMgr struct {
	// 连接索引到玩家的映射
	playerConnIdxMap map[int64]*Player
	ConnIdxLock      sync.RWMutex

	// 用户和平台类型到玩家的映射
	userAndPlatformMap  map[string]*Player
	UserAndPlatformLock sync.RWMutex

	// 玩家ID到玩家的映射
	accountIDMap  map[uint64]*Player
	AccountIDLock sync.RWMutex

	ConnIdx int64
}

var (
	PlayerMgr *playerMgr
)

func init() {
	fmt.Println("playerMgr init")
	PlayerMgr = &playerMgr{}
	PlayerMgr.Init()
}

func (pm *playerMgr) Init() {
	pm.playerConnIdxMap = make(map[int64]*Player)
	pm.userAndPlatformMap = make(map[string]*Player)
	pm.accountIDMap = make(map[uint64]*Player)

	pm.ConnIdx = time.Now().Unix()
}

func (pm *playerMgr) CreatePlayer(conn gnet.Conn) *Player {
	connIdx := pm.ConnIdx
	pm.ConnIdx++
	player := &Player{}
	player.Init(conn, connIdx)

	pm.AddPlayerByConnIdx(connIdx, player)

	return player
}

func (pm *playerMgr) DeletePlayer(player *Player) {
	pm.ConnIdxLock.Lock()
	delete(pm.playerConnIdxMap, player.ConnIdx)
	pm.ConnIdxLock.Unlock()

	pm.UserAndPlatformLock.Lock()
	delete(pm.userAndPlatformMap, player.User+"_"+player.Platform)
	pm.UserAndPlatformLock.Unlock()

	pm.AccountIDLock.Lock()
	delete(pm.accountIDMap, player.AccountID)
	pm.AccountIDLock.Unlock()
}

func (pm *playerMgr) GetPlayerByConnIdx(conn_idx int64) *Player {
	pm.ConnIdxLock.RLock()
	defer pm.ConnIdxLock.RUnlock()

	return pm.playerConnIdxMap[conn_idx]
}

func (pm *playerMgr) GetPlayerByUserAndPlatform(user string, plat_type string) *Player {
	pm.UserAndPlatformLock.RLock()
	defer pm.UserAndPlatformLock.RUnlock()

	return pm.userAndPlatformMap[user+"_"+plat_type]
}

func (pm *playerMgr) GetPlayerByAccountID(account_id uint64) *Player {
	pm.AccountIDLock.RLock()
	defer pm.AccountIDLock.RUnlock()

	return pm.accountIDMap[account_id]
}

func (pm *playerMgr) AddPlayerByUserAndPlatform(player *Player) {
	pm.UserAndPlatformLock.Lock()
	defer pm.UserAndPlatformLock.Unlock()

	pm.userAndPlatformMap[player.User+"_"+player.Platform] = player
}

func (pm *playerMgr) AddPlayerByAccountID(player *Player) {
	pm.AccountIDLock.Lock()
	defer pm.AccountIDLock.Unlock()

	pm.accountIDMap[player.AccountID] = player
}

func (pm *playerMgr) AddPlayerByConnIdx(ConnIdx int64, player *Player) {
	pm.ConnIdxLock.Lock()
	defer pm.ConnIdxLock.Unlock()

	pm.playerConnIdxMap[ConnIdx] = player
}

func (pm *playerMgr) Tick1S() {
	// fmt.Println("playerMgr Tick1S")
}

func (pm *playerMgr) LoginOut(player *Player, reason uint32) {
	if player == nil {
		return
	}

	pm.DeletePlayer(player)

	if player.GetConn() != nil {
		player.GetConn().Close()
	}
}

func (pm *playerMgr) KickAllPlayer(reason uint32) {
	allPlayer := make([]*Player, len(pm.playerConnIdxMap))

	pm.ConnIdxLock.Lock()
	for _, player := range pm.playerConnIdxMap {
		allPlayer = append(allPlayer, player)
	}
	pm.ConnIdxLock.Unlock()

	for _, player := range allPlayer {
		pm.LoginOut(player, reason)
	}
}
