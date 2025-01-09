package model

type ActiveMap struct {
	// 节点ID
	NodeID uint32
	// 子节点执行类型
	NodeType uint32
	// 子节点列表
	NodeList []uint32
	// 执行动作
	Action func(agent *Agent)
	// 执行状态
	NodeState uint32
}

const (
	// 节点类型
	ACTIVE_NODE_TYPE_UNKNOW = 0
	ACTIVE_NODE_TYPE_ORDER  = 1
	ACTIVE_NODE_TYPE_RANDOM = 2
	ACTIVE_NODE_TYPE_SUB    = 3
)

const (
	ACTIVE_TYPE_SCENE = 1
)
