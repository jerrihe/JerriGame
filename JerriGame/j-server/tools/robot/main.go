package main

import (
	"fmt"
	// snow "jserver/src/common/snowflake"
	mainserver "robot/mainserver"
	model "robot/model"
)

func main() {
	fmt.Println("start robot!")
	// snow.InitSnowFlake()
	model.AgentMgrIn.Init()
	model.InitActive()
	model.InitScene()

	mainserver.StartServer()
}
