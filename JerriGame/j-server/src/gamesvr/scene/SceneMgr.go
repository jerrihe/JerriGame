package scene

import (
	"fmt"
	model "gamesvr/model"
	"sync"
)

type SceneMgr struct {
	Scenes     map[uint32]map[uint32]*Scene
	SceneSLock sync.RWMutex
}

var SceneMgrIn *SceneMgr

func init() {
	SceneMgrIn = &SceneMgr{}
	SceneMgrIn.Init()
}

func (sm *SceneMgr) Init() {
	sm.Scenes = make(map[uint32]map[uint32]*Scene)
}

func (sm *SceneMgr) GetScene(map_id uint32, scene_id uint32) *Scene {
	sm.SceneSLock.RLock()
	defer sm.SceneSLock.RUnlock()

	return sm.Scenes[map_id][scene_id]
}

func (sm *SceneMgr) AddScene(map_id uint32, scene_id uint32) {
	sm.SceneSLock.Lock()
	defer sm.SceneSLock.Unlock()

	scene := &Scene{}
	scene.Init(scene_id, map_id)

	if sm.Scenes[scene.MapID] == nil {
		sm.Scenes[scene.MapID] = make(map[uint32]*Scene)
	}
	sm.Scenes[scene.MapID][scene.SceneID] = scene
	fmt.Println("AddScene map_id: ", map_id, " scene_id: ", scene_id)
}

func (sm *SceneMgr) DelScene(map_id uint32, scene_id uint32) {
	sm.SceneSLock.Lock()
	defer sm.SceneSLock.Unlock()

	delete(sm.Scenes[map_id], scene_id)
	fmt.Println("DelScene map_id: ", map_id, " scene_id: ", scene_id)
}

func (sm *SceneMgr) GetSceneByMapID(map_id uint32) map[uint32]*Scene {
	sm.SceneSLock.RLock()
	defer sm.SceneSLock.RUnlock()

	return sm.Scenes[map_id]
}

func (sm *SceneMgr) GetSceneByMapIDSceneID(map_id uint32, scene_id uint32) *Scene {
	sm.SceneSLock.RLock()
	defer sm.SceneSLock.RUnlock()

	return sm.Scenes[map_id][scene_id]
}

func (sm *SceneMgr) EnterPlayer(player *model.Player, map_id uint32, scene_id uint32) bool {

	var scene *Scene
	scene = sm.GetScene(map_id, scene_id)
	if scene == nil {
		fmt.Println("EnterPlayer scene is nil")

		// 场景不存在 创建场景
		sm.AddScene(map_id, scene_id)
		scene = sm.GetScene(map_id, scene_id)
		if scene == nil {
			return false
		}
	}
	return scene.EnterPlayer(player)
}

func (sm *SceneMgr) LeavePlayer(player *model.Player, map_id uint32, scene_id uint32) bool {
	scene := sm.GetScene(map_id, scene_id)
	if scene == nil {
		return false
	}

	return scene.LeavePlayer(player)
}
