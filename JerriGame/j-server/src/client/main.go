package main

import (
	"fmt"
	"time"

	app "client/app"

	pool "client/workpool"
)

func main() {
	pool.InitWorkPool(5)
	app.InitApp()
	// 服务端地址
	serverAddress := "192.168.31.97:9015"

	// 运行客户端
	min_robot := 1
	max_robot := 2

	for i := min_robot; i < max_robot; i++ {
		go app.RunAgentClient(serverAddress, fmt.Sprint(i), "1")
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
