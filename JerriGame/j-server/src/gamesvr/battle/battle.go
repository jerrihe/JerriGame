package battle

type Record struct {
	// 出手方
	order int
	// 伤害值
	damage int
}

type War struct {
	blue BattleInfo
	red  BattleInfo

	// 出手顺序
	next_order int
	// 战斗记录
	record []Record
}

func (w *War) Init(blue BattleInfo, red BattleInfo) {
	w.blue = blue
	w.red = red
}

func (w *War) GetBlue() BattleInfo {
	return w.blue
}

func (w *War) GetRed() BattleInfo {
	return w.red
}

func (w *War) Start() {
	// 初始化战斗出手顺序
	w.next_order = PLAYER
}

func (w *War) Step() {
	// 出手
	var attacker, defender *BattleInfo
	if w.next_order == PLAYER {
		attacker = &w.blue
		defender = &w.red
	} else {
		attacker = &w.red
		defender = &w.blue
	}

	// 计算伤害
	damage := attacker.Attack
	defender.HP -= damage

	// 记录战斗
	w.record = append(w.record, Record{w.next_order, damage})

	// 切换出手方
	w.next_order = 1 - w.next_order
}
