package main

import (
	// "gamesvr/game"
	"gamesvr/mainserver"

	app "gamesvr/app"
	dao "gamesvr/dao"
	// snow "jserver/src/common/snowflake"
)

func main() {
	// game.PlayerMgr.Init()
	app.InitApp()
	dao.InitDaoMgr()
	// snow.InitSnowFlake()

	mainserver.StartServer()
}
