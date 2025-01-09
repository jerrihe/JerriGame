package scene

import (
	"container/heap"
	"fmt"
	"math"
)

// 节点结构
type Node struct {
	X, Y     int
	Cost     float64
	Priority float64
	Parent   *Node
}

// 优先队列实现
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq *PriorityQueue) Swap(i, j int) { (*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// 启发式函数类型
type HeuristicFunc func(x1, y1, x2, y2 int) float64
type HeuristicDistanceFunc func(p1, p2 Pos) float64

// 启发式函数定义
// 曼哈顿距离
func Manhattan(x1, y1, x2, y2 int) float64 {
	return math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2))
}

// 计算曼哈顿距离
func ManhattanDistance(p1, p2 Pos) float64 {
	return math.Abs(p2.X-p1.X) + math.Abs(p2.Y-p1.Y)
}

// 欧几里得距离
func Euclidean(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2))
}

// 计算欧几里得距离
func EuclideanDistance(p1, p2 Pos) float64 {
	return math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2))
}

// 切比雪夫距离
func Chebyshev(x1, y1, x2, y2 int) float64 {
	return math.Max(math.Abs(float64(x1-x2)), math.Abs(float64(y1-y2)))
}

// 计算切比雪夫距离
func ChebyshevDistance(p1, p2 Pos) float64 {
	return math.Max(math.Abs(p2.X-p1.X), math.Abs(p2.Y-p1.Y))
}

// 无启发式函数
func ZeroHeuristic(x1, y1, x2, y2 int) float64 {
	return 0
}

// 检查是否在地图范围内
func isValid(x, y int) bool {
	return x >= 0 && x < WIDTH && y >= 0 && y < HEIGHT && grid[y][x] == 0
}

// A* 算法
func AStar(startX, startY, endX, endY int, heuristic HeuristicFunc) []*Node {
	// 支持对角线移动
	directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	openSet := &PriorityQueue{}
	heap.Init(openSet)
	startNode := &Node{X: startX, Y: startY, Cost: 0, Priority: heuristic(startX, startY, endX, endY)}
	heap.Push(openSet, startNode)
	closedSet := make(map[[2]int]bool)

	for openSet.Len() > 0 {
		// 取出开放集中的第一个节点
		current := heap.Pop(openSet).(*Node)
		// 到达终点
		if current.X == endX && current.Y == endY {
			var path []*Node
			for current != nil {
				path = append([]*Node{current}, path...)
				current = current.Parent
			}
			return path
		}

		closedSet[[2]int{current.X, current.Y}] = true
		// 遍历当前节点的邻居
		for _, dir := range directions {
			neighborX, neighborY := current.X+dir[0], current.Y+dir[1]
			if !isValid(neighborX, neighborY) || closedSet[[2]int{neighborX, neighborY}] {
				continue
			}
			cost := current.Cost + 1
			if dir[0] != 0 && dir[1] != 0 {
				cost = current.Cost + math.Sqrt(2)
			}
			neighbor := &Node{
				X:        neighborX,
				Y:        neighborY,
				Cost:     cost,
				Priority: cost + heuristic(neighborX, neighborY, endX, endY),
				Parent:   current,
			}
			heap.Push(openSet, neighbor)
		}
	}
	return nil
}

// 转换实际坐标到网格坐标
func ToGridCoords(xActual, yActual float64) (int, int) {
	return int(xActual / CELL), int(yActual / CELL)
}

// 转换网格坐标到实际坐标
func ToActualCoords(xGrid, yGrid int) (float64, float64) {
	return float64(xGrid) * CELL, float64(yGrid) * CELL
}

func GetPathWithDistances(path []Pos) ([]Pos, float64) {
	if len(path) < 2 {
		// 如果路径点少于 2 个，则没有距离可计算
		return path, 0
	}

	totalDistance := 0.0

	// 遍历路径点，计算到下一个点的距离
	for i := 0; i < len(path)-1; i++ {
		current := path[i]
		next := path[i+1]

		// 计算欧几里得距离
		distance := EuclideanDistance(current, next)

		// 更新当前点的距离到下一个点
		path[i].DistanceToNext = distance

		// 累加总距离
		totalDistance += distance
	}

	// 最后一个点到下一个点的距离为 0
	path[len(path)-1].DistanceToNext = 0

	return path, totalDistance
}

func GetPath(startXActual, startYActual, endXActual, endYActual float64) ([]Pos, float64) {
	// 启发式函数
	heuristic := Manhattan
	// 启发式函数距离
	heuristicDistance := ManhattanDistance

	startX, startY := ToGridCoords(startXActual, startYActual)
	endX, endY := ToGridCoords(endXActual, endYActual)
	if !isValid(startX, startY) {
		fmt.Println("起点无效（被障碍物占据或超出地图范围）！")
		return nil, 0
	}

	if !isValid(endX, endY) {
		fmt.Println("终点无效（被障碍物占据或超出地图范围）！")
		return nil, 0
	}

	path := AStar(startX, startY, endX, endY, heuristic)

	if path != nil {
		// 初始化实际路径和总距离
		var actualPath []Pos = make([]Pos, len(path))
		totalDistance := 0.0
		for i := 0; i < len(path); i++ {
			x, y := ToActualCoords(path[i].X, path[i].Y)
			actualPath = append(actualPath, Pos{X: x, Y: y})

			if i > 0 {
				// 计算欧几里得距离
				actualPath[i-1].DistanceToNext = heuristicDistance(actualPath[i-1], actualPath[i])
				totalDistance += actualPath[i-1].DistanceToNext
			}
		}
		return actualPath, totalDistance
	}
	return nil, 0
}
