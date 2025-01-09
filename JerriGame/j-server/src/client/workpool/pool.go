package workpool

import (
	"fmt"
	"sync"
)

// Task 是一个任务接口
type Task struct {
	ID   int
	Task func()
}

// WorkerPool 是协程池
type WorkerPool struct {
	workerCount int            // 最大并发数
	tasks       chan Task      // 任务队列
	wg          sync.WaitGroup // 用于等待所有任务完成
}

var Worker *WorkerPool

// InitWorkPool 初始化协程池
func InitWorkPool(workerCount int) {
	Worker = NewWorkerPool(workerCount)
	Worker.Run()
}

// NewWorkerPool 创建一个新的协程池
func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		tasks:       make(chan Task),
	}
}

// Run 启动协程池
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.workerCount; i++ {
		fmt.Printf("Worker %d started\n", i)
		go wp.worker(i)
	}
}

// worker 是每个工作 Goroutine 的逻辑
func (wp *WorkerPool) worker(workerID int) {
	for task := range wp.tasks {
		task.Task() // 执行任务
		fmt.Printf("Worker %d processing task %d\n", workerID, task.ID)
		wp.wg.Done() // 通知任务完成
	}
}

// AddTask 添加任务到任务队列
func (wp *WorkerPool) AddTask(task Task) {
	wp.wg.Add(1) // 增加任务计数
	wp.tasks <- task
}

// Wait 等待所有任务完成
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.tasks) // 关闭任务队列，通知所有 worker 退出
}
