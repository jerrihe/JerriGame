package battle

const (
	HP = iota
	ATTACK
)

const (
	PLAYER = iota
	ENEMY
)

type BattleInfo struct {
	HP     int
	Attack int
}

func (b *BattleInfo) Init(battle_role_type int) {
	switch battle_role_type {
	case PLAYER:
		b.HP = 100
		b.Attack = 10
	case ENEMY:
		b.HP = 100
		b.Attack = 5
	default:
		return
	}
}
