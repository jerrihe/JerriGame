package app

import (
	"fmt"
	"gamesvr/game"
	message "jserver/src/common/message"
	snow "jserver/src/common/snowflake"
	"strconv"
	"time"

	conf "gamesvr/config"
	dao "gamesvr/dao"
	model "gamesvr/model"
	ServerMgr "gamesvr/server"

	"jserver/src/protocol/cs"
	"jserver/src/protocol/ss"

	"google.golang.org/protobuf/proto"
)

type HttpBody struct {
	Code int `json:"code"`
}

func init() {
	fmt.Println("app init")
	game.RegisterHandle(cs.CS_CMD_LOGIN_REQ, HanleCsCmdLoginReq)
	game.RegisterHandle(cs.CS_CMD_CREATE_ACCOUNT_REQ, HandleCsCmdCreateAccountReq)
	game.RegisterHandle(cs.CS_CMD_LOGIN_OUT_REQ, HandleCsCmdLoginOutReq)

	game.RegisterHandleSS(ss.SS_CMD_KICK_ACCOUNT_REQ, HandleSsCmdKickAccountReq)
	game.RegisterHandleSS(ss.SS_CMD_KICK_ACCOUNT_RES, HandleKickRoleRes)
}

// func checkToken(user string, player_type string) bool {
// 	// 平台校验token

// 	targetUrl := "http://192.168.31.97:8000/api/loginsvr/v1/checkAccess"

// 	params := map[string]interface{}{
// 		"user":     user,
// 		"platType": player_type,
// 	}
// 	jsonData, err := json.Marshal(params)
// 	if err != nil {
// 		fmt.Println("Error marshalling data:", err)
// 		return false
// 	}

// 	res, err := http.Post(targetUrl, "application/json", bytes.NewReader(jsonData))
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return false
// 	}

// 	defer res.Body.Close()

// 	body, _ := io.ReadAll(res.Body)
// 	fmt.Println(string(body))

// 	var httpBody HttpBody
// 	if err := json.Unmarshal(body, &httpBody); err != nil {
// 		fmt.Println("Error unmarshalling data:", err)
// 		return false
// 	}

// 	if httpBody.Code != 0 {
// 		fmt.Println("check token failed")
// 		return false
// 	}

// 	return true
// }

func AccountLogin(player *model.Player, account uint64) {
	player.SetAccountID(account)

	model.PlayerMgr.AddPlayerByAccountID(player)
}

func CheckOnlineLogin(player *model.Player, user string, platform string) bool {
	// 查看是否在其他服登录
	if exists, err := dao.RedisMgrInstance.HExists("login", user); err == nil && exists {
		fmt.Println("user is login")
		// 去其他服踢掉账号
		or_server_id, err := dao.RedisMgrInstance.HGet("login", user)
		fmt.Println("or_server_id:", or_server_id, "server_id:", conf.Conf.ServerInfo.ServerID)
		if err != nil {
			fmt.Println("HGet failed. Error: ", err)
		} else {
			num64, err := strconv.ParseInt(or_server_id, 10, 32)
			if err != nil {
				fmt.Println("ParseInt failed. Error: ", err)
			}

			// 不是本服才踢
			if num64 == int64(conf.Conf.ServerInfo.ServerID) {
				return false
			}

			dao.RedisMgrInstance.HSet("login", user, conf.Conf.ServerInfo.ServerID)
			fmt.Println("other server:", or_server_id)

			// 发送踢人消息
			var kickReq ss.SsCmdKickAccountReq
			kickReq.User = user
			kickReq.Platform = platform
			kickReq.ConnIdx = uint64(player.GetConnIdx())

			message.SendToRouter(ServerMgr.ServerMgr.GetServerConn(), ss.SS_CMD_KICK_ACCOUNT_REQ, conf.Conf.ServerInfo.ServerID, int32(num64), &kickReq)
		}
		return true
	} else {
		fmt.Println("HExists failed. Error: ", err)
		dao.RedisMgrInstance.HSet("login", user, conf.Conf.ServerInfo.ServerID)
	}
	return false
}

func HanleCsCmdLoginReq(player *model.Player, head *cs.CsHead, body *proto.Message) {
	fmt.Println("HanleLogin")

	if player == nil || body == nil {
		fmt.Println("player is nil")
		return
	}

	loginReq, ok := (*body).(*cs.CsCmdLoginReq)
	if !ok {
		fmt.Println("loginReq is nil")
		return
	}

	// 检查是否已经登录
	if player.IsLogin() {
		fmt.Println("player is login")
		player.SendErrCode(cs.CS_CMD_LOGIN_RES, head.Seq, cs.ERR_CODE_LOGIN_ONLINE)
		return
	}

	// 防止客户端伪造
	if (loginReq.User != "" && player.GetUser() != "" && loginReq.User != player.GetUser()) || (loginReq.Platform != "" && player.GetPlatform() != "" && loginReq.Platform != player.GetPlatform()) {
		fmt.Println("user or platform is not match", loginReq.User, loginReq.Platform, player.GetUser(), player.GetPlatform())
		player.SendErrCode(cs.CS_CMD_LOGIN_RES, head.Seq, cs.ERR_CODE_LOGIN_INVALID_ACCOUNT)
		return
	}

	player.InitLogin(loginReq.User, loginReq.Platform)

	// 平台校验token
	// if !checkToken(loginReq.User, loginReq.Platform) {
	// 	fmt.Println("check token failed")
	// 	return
	// }

	// 检查数据库是否账号
	var accDataList []dao.DBAccData
	// 查询玩家账号条件
	accWhere := fmt.Sprintf("UserId = \"%s\" and Platform = \"%s\" ", loginReq.User, loginReq.Platform)

	// 获取账号数据
	if err := dao.GetTableData(dao.GameDB, dao.AccountTable, 0, &accDataList, accWhere); err != nil {
		fmt.Println("GetTableData failed. Error: ", err)
		player.SendErrCode(cs.CS_CMD_LOGIN_RES, head.Seq, cs.ERR_CODE_LOGIN_INVALID_ACCOUNT)
		return
	}

	// 如果没有账号则提醒创建账号
	if len(accDataList) == 0 {
		player.SendErrCode(cs.CS_CMD_LOGIN_RES, head.Seq, cs.ERR_CODE_LOGIN_NEED_CREATE)
		return
	}

	// 查看是否在其他服登录
	if CheckOnlineLogin(player, loginReq.User, loginReq.Platform) {
		return
	}

	// 检查同服登录
	if oldPlayer := model.PlayerMgr.GetPlayerByUserAndPlatform(loginReq.User, loginReq.Platform); oldPlayer != nil {
		// 发送踢人消息
		var kickntf cs.CsCmdNtfKickAccount
		kickntf.Reason = 1
		oldPlayer.SendMsg(cs.CS_CMD_NTF_KICK_ACCOUNT, 0, &kickntf)

		// 删除链接 由对方删除链接
		model.PlayerMgr.LoginOut(oldPlayer, 1)
	}

	// 返回账号信息
	var loginRes cs.CsCmdLoginRes
	loginRes.AccountId = accDataList[0].AccId

	model.PlayerMgr.AddPlayerByUserAndPlatform(player)

	// 处理账号登录
	AccountLogin(player, accDataList[0].AccId)

	player.SendMsg(cs.CS_CMD_LOGIN_RES, head.Seq, &loginRes)
}

func HandleCsCmdCreateAccountReq(player *model.Player, head *cs.CsHead, body *proto.Message) {
	if player == nil {
		fmt.Println("player is nil")
		return
	}

	var req = (*body).(*cs.CsCmdCreateAccountReq)
	fmt.Printf("CreateRoleInfo: %#v\n", req)

	// 查询redis中有没有账号

	var accDataList []dao.DBAccData
	// 查询玩家账号条件
	accWhere := fmt.Sprintf("UserId = \"%s\" and Platform = \"%s\" ", req.User, req.Platform)
	// 获取账号数据
	if err := dao.GetTableData(dao.GameDB, dao.AccountTable, 0, &accDataList, accWhere); err != nil {
		fmt.Println("GetTableData failed. Error: ", err)
		player.SendErrCode(cs.CS_CMD_CREATE_ACCOUNT_RES, head.Seq, cs.ERR_CODE_LOGIN_INVALID_ACCOUNT)
		return
	}

	fmt.Println("accDataList:", accDataList)

	// 如果没有账号则提醒创建账号
	if len(accDataList) > 0 {
		player.SendErrCode(cs.CS_CMD_CREATE_ACCOUNT_RES, head.Seq, cs.ERR_CODE_CREATE_EXITE_ACCOUNT)
		return
	}

	// 创建账号
	account_id := snow.SF.NextID()

	// 插入账号数据
	accData := dao.DBAccData{
		UserId:   req.User,
		Platform: req.Platform,
		AccId:    uint64(account_id),
		CreateAt: uint64(time.Now().Unix()),
	}

	var res cs.CsCmdCreateAccountRes
	if err := dao.InsertData(dao.GameDB, dao.AccountTable, 0, &dao.DBAccData{}, &accData); err != nil {
		fmt.Println("InsertTableData failed. Error: ", err)
		player.SendErrCode(cs.CS_CMD_CREATE_ACCOUNT_RES, head.Seq, cs.ERR_CODE_CREATE_ACCOUNT_FAILED)
		return
	} else {
		res.AccountId = uint64(account_id)
		player.SetAccountID(uint64(account_id))
	}

	player.SendMsg(cs.CS_CMD_CREATE_ACCOUNT_RES, head.Seq, &res)
}

func HandleCsCmdLoginOutReq(player *model.Player, head *cs.CsHead, body *proto.Message) {
	if player == nil || body == nil {
		fmt.Println("player is nil")
		return
	}

	var req = (*body).(*cs.CsCmdLoginOutReq)
	fmt.Printf("LoginOutInfo: %#v\n", req)
	var res cs.CsCmdLoginOutRes
	player.SendMsg(cs.CS_CMD_LOGIN_OUT_RES, head.Seq, &res)

	// 登出
	model.PlayerMgr.LoginOut(player, 0)

	// 删除redis中账号
	dao.RedisMgrInstance.HDel("login", player.GetUser())
}

func HandleSsCmdKickAccountReq(head *ss.SsHead, body *proto.Message) {

	req := (*body).(*ss.SsCmdKickAccountReq)

	fmt.Printf("KickRoleInfo:%#v\n", req)

	// 踢人
	if req.User == "" {
		fmt.Println("kick user is nil")
		return
	}

	if player := model.PlayerMgr.GetPlayerByUserAndPlatform(req.User, req.Platform); player != nil {
		// 发送踢人消息
		var kickntf cs.CsCmdNtfKickAccount
		kickntf.Reason = 1
		player.SendMsg(cs.CS_CMD_NTF_KICK_ACCOUNT, 0, &kickntf)

		// 删除链接 由对方删除链接
		// player.GetConn().Close()

		// 删除账号
		model.PlayerMgr.LoginOut(player, 1)
		fmt.Println("kick user success", req.User, req.Platform)
	}

	// 返回踢人结果
	var res ss.SsCmdKickAccountRes
	res.User = req.User
	res.Platform = req.Platform
	res.ConnIdx = req.ConnIdx

	message.SendToRouter(ServerMgr.ServerMgr.GetServerConn(), ss.SS_CMD_KICK_ACCOUNT_RES, conf.Conf.ServerInfo.ServerID, head.ServerId, &res)
}

func HandleKickRoleRes(head *ss.SsHead, body *proto.Message) {

	res := (*body).(*ss.SsCmdKickAccountRes)

	fmt.Printf("KickRoleRes:%#v\n", res)

	//登录
	if res.User == "" {
		fmt.Println("kick user is nil")
	}

	player := model.PlayerMgr.GetPlayerByConnIdx(int64(res.ConnIdx))
	if player == nil {
		fmt.Println("player is nil")
	}

	// 登录
	// 检查数据库是否账号
	var accDataList []dao.DBAccData
	// 查询玩家账号条件
	accWhere := fmt.Sprintf("UserId = \"%s\" and Platform = \"%s\" ", res.User, res.Platform)

	// 获取账号数据
	if err := dao.GetTableData(dao.GameDB, dao.AccountTable, 0, &accDataList, accWhere); err != nil {
		fmt.Println("GetTableData failed. Error: ", err)
		return
	}

	fmt.Println("accDataList:", accDataList)

	// 返回账号信息
	var loginRes cs.CsCmdLoginOutRes
	loginRes.AccountId = accDataList[0].AccId
	player.SetAccountID(accDataList[0].AccId)

	model.PlayerMgr.AddPlayerByAccountID(player)

	player.SendMsg(cs.CS_CMD_LOGIN_RES, 0, &loginRes)
}
