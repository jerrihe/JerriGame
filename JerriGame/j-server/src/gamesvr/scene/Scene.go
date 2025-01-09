package scene

import (
	"fmt"
	model "gamesvr/model"
	"jserver/src/protocol/cs"
	"sync"
)

type Scene struct {
	// 场景ID
	SceneID uint32
	// mapID
	MapID uint32
	// 场景进入玩家
	EnterPlayers     map[uint64]*model.Player
	EnterPlayersLock sync.RWMutex

	GridPlayers map[int]map[uint64]*model.Player
}

func (s *Scene) Init(scene_id uint32, map_id uint32) {
	s.SceneID = scene_id
	s.MapID = map_id
	s.EnterPlayers = make(map[uint64]*model.Player)
	s.GridPlayers = make(map[int]map[uint64]*model.Player)
	fmt.Println("Scene Init scene_id: ", scene_id, " map_id: ", map_id)
}

func (s *Scene) EnterPlayer(player *model.Player) bool {
	s.EnterPlayersLock.Lock()
	defer s.EnterPlayersLock.Unlock()

	if s.EnterPlayers[player.GetAccountID()] != nil {
		return false
	}

	s.EnterPlayers[player.GetAccountID()] = player

	player.InitSceneInfo(s.SceneID, s.MapID, 0, 0, 0)

	fmt.Println("Scene EnterPlayer account_id", player.GetAccountID())
	return true
}

func (s *Scene) LeavePlayer(player *model.Player) bool {
	s.EnterPlayersLock.Lock()
	defer s.EnterPlayersLock.Unlock()

	account_id := player.GetAccountID()
	if s.EnterPlayers[account_id] == nil {
		return false
	}

	// 处理视野
	s.ChangeGrid(account_id, -1, -1)

	// 离开场景
	delete(s.EnterPlayers, account_id)
	player.ClearSceneInfo()
	return true
}

// 改变格子 用于测试 x,y为格子坐标 1格子长宽=5米 Grid : GridMaxRow*GridMaxCol
func (s *Scene) ChangeGrid(account_id uint64, x float32, y float32) bool {
	s.EnterPlayersLock.Lock()
	defer s.EnterPlayersLock.Unlock()

	player := s.EnterPlayers[account_id]
	if player == nil {
		return false
	}

	old_x := player.SceneInfo.Pos.X
	old_y := player.SceneInfo.Pos.Y

	if x == old_x && y == old_y {
		return false
	}

	var newGrid int32 = -1
	if x >= 0 && y >= 0 && x < float32(GridHeight*GridMaxRow) && y < float32(GridHeight*GridMaxCol) {
		newGrid = int32(x/float32(GridHeight)) + int32(y/float32(GridHeight*GridMaxCol))
	} else {
		// 离开场景
		fmt.Println("Player account_id", account_id, " x", x, " y", y, " leave scene")
	}

	// 更新位置
	player.SceneInfo.Pos.X = x
	player.SceneInfo.Pos.Y = y
	player.SceneInfo.OldPos.X = old_x
	player.SceneInfo.OldPos.Y = old_y
	player.SceneInfo.OldGrid = player.SceneInfo.GridID
	player.SceneInfo.GridID = int32(newGrid)

	// 通知视野中的人我移动了
	s.NotifySeePlayerMove(player)

	delete(s.GridPlayers, int(player.SceneInfo.OldGrid))

	fmt.Println("Player account_id", account_id, " old_x", old_x, " old_y", old_y, " x", x, " y", y)

	oldGrig := player.SceneInfo.OldGrid

	if oldGrig != newGrid {
		fmt.Println("Player account_id", account_id, " oldGrid", oldGrig, " newGrid", newGrid)
		// 更新玩家到新格子
		if _, ok := s.GridPlayers[int(newGrid)]; !ok {
			s.GridPlayers[int(newGrid)] = make(map[uint64]*model.Player)
		}
		s.GridPlayers[int(newGrid)][account_id] = player

		// 从旧格子删除
		if _, ok := s.GridPlayers[int(oldGrig)]; ok {
			delete(s.GridPlayers[int(oldGrig)], account_id)
		}

		// TODO: 广播给周围玩家
		oldGrids := s.GetNeighbors(int(oldGrig))
		newGrids := s.GetNeighbors(int(newGrid))

		fmt.Println("Player account_id: ", account_id, " oldGrids: ", oldGrids, " newGrids: ", newGrids)

		levelPlayers := make([]*model.Player, 0)
		for k := range oldGrids {
			if _, ok := newGrids[k]; !ok {
				// 离开的格子
				if playerList, ok := s.GridPlayers[k]; ok {
					for _, p := range playerList {
						if p.GetAccountID() != account_id {
							levelPlayers = append(levelPlayers, p)
						}
					}
				}
			}
		}
		enterPlayers := make([]*model.Player, 0)
		for k := range newGrids {
			if _, ok := oldGrids[k]; !ok {
				// 进入的格子
				if playerList, ok := s.GridPlayers[k]; ok {
					for _, p := range playerList {
						if p.GetAccountID() != account_id {
							enterPlayers = append(enterPlayers, p)
						}
					}
				}
			}
		}
		fmt.Println("Player account_id: ", account_id, " leaveGrids: ", levelPlayers, " enterGrids: ", enterPlayers)
		if len(levelPlayers) > 0 {
			// TODO: 离开周围玩家视野并广播给周围玩家
			s.DealWithLeaveSeePlayer(player, levelPlayers)
		}

		if len(enterPlayers) > 0 {
			// TODO: 进入周围玩家视野并广播给周围玩家
			s.DealWithEnterSeePlayer(player, enterPlayers)
		}
	}
	return true
}

// 获取给定索引的邻居索引列表
func (s *Scene) GetNeighbors(gridID int) map[int]bool {
	neighbors := make(map[int]bool, 0)

	if gridID < 0 {
		return neighbors
	}

	row := gridID / GridMaxRow
	col := gridID % GridMaxRow

	// 遍历相对的行和列位置
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			// 跳过中心点本身
			// if dr == 0 && dc == 0 {
			// 	continue
			// }

			// 计算邻居的行和列
			neighborRow := row + dr
			neighborCol := col + dc

			// 检查是否在边界内
			if neighborRow >= 0 && neighborRow < GridMaxRow && neighborCol >= 0 && neighborCol < GridMaxRow {
				// 计算邻居的索引
				neighborIndex := neighborRow*GridMaxRow + neighborCol
				// neighbors = append(neighbors, neighborIndex)
				neighbors[neighborIndex] = true
			}
		}
	}

	return neighbors
}

// 处理离开视野
func (s *Scene) DealWithLeaveSeePlayer(player *model.Player, leavePlayes []*model.Player) {
	if player == nil || len(leavePlayes) == 0 {
		return
	}

	var account_id = player.GetAccountID()
	var otherNtf cs.CsCmdSceneNtfLeave
	otherNtf.AccountId = append(otherNtf.AccountId, account_id)

	var ntf cs.CsCmdSceneNtfLeave
	for _, p := range leavePlayes {
		// 我离开了p的视野
		p.DelSeePlayer(account_id)
		// TODO: 通知p我离开了
		p.SendMsg(cs.CS_CMD_SCENE_NTF_LEAVE, 0, &otherNtf)

		// p离开了我的视野
		player.DelSeePlayer(p.GetAccountID())
		ntf.AccountId = append(ntf.AccountId, p.GetAccountID())
	}
	// TODO: 通知我离开我视野的玩家列表
	player.SendMsg(cs.CS_CMD_SCENE_NTF_LEAVE, 0, &ntf)
}

// 处理进入视野
func (s *Scene) DealWithEnterSeePlayer(player *model.Player, enterPlayes []*model.Player) {
	if player == nil || len(enterPlayes) == 0 {
		return
	}

	var account_id = player.GetAccountID()
	var ntf cs.CsCmdSceneNtfEnter
	var otherNtf cs.CsCmdSceneNtfEnter
	otherNtf.AccountId = append(otherNtf.AccountId, account_id)

	for _, p := range enterPlayes {
		// 我进入了p的视野
		p.AddSeePlayer(account_id)
		// TODO: 通知p我进入了
		p.SendMsg(cs.CS_CMD_SCENE_NTF_ENTER, 0, &otherNtf)
		// p进入了我的视野
		player.AddSeePlayer(p.GetAccountID())
		ntf.AccountId = append(ntf.AccountId, p.GetAccountID())
	}
	// TODO: 通知我进入我视野的玩家列表
	player.SendMsg(cs.CS_CMD_SCENE_NTF_ENTER, 0, &ntf)
}

// 通知其余人移动
func (s *Scene) NotifyMove(player *model.Player, x float32, y float32) bool {
	s.EnterPlayersLock.RLock()
	defer s.EnterPlayersLock.RUnlock()

	account_id := player.GetAccountID()
	if s.EnterPlayers[account_id] == nil {
		return false
	}

	// 开始移动
	s.ChangeGrid(account_id, x, y)

	return true
}

// 通知我视野中的人我的移动
func (s *Scene) NotifySeePlayerMove(player *model.Player) {
	var ntf cs.CsCmdSceneNtfMove
	ntf.AccountId = player.AccountID
	ntf.X = uint32(player.SceneInfo.Pos.X * 100)
	ntf.Y = uint32(player.SceneInfo.Pos.Y * 100)
	ntf.OldX = uint32(player.SceneInfo.OldPos.X * 100)
	ntf.OldY = uint32(player.SceneInfo.OldPos.Y * 100)

	for k := range player.SceneInfo.SeePlayers {
		if p := s.EnterPlayers[k]; p != nil {
			p.SendMsg(cs.CS_CMD_SCENE_NTF_MOVE, 0, &ntf)
		}
	}
}
