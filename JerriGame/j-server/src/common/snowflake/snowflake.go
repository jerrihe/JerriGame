package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// 雪花算法生成器结构
type Snowflake struct {
	mu           sync.Mutex // 保证并发安全
	timestamp    int64      // 上次生成 ID 的时间戳
	workerID     int64      // 当前机器 ID
	sequence     int64      // 当前毫秒内的序列号
	epoch        int64      // 起始时间戳（自定义）
	workerIDBits int64      // 机器 ID 所占的位数
	sequenceBits int64      // 序列号所占的位数
	maxWorkerID  int64      // 最大机器 ID
	maxSequence  int64      // 最大序列号
	timeShift    int64      // 时间戳左移位数
	workerShift  int64      // 机器 ID 左移位数
}

var SF *Snowflake

func init() {
	fmt.Printf("SnowFlake init\n")
	InitSnowFlake()
}

func InitSnowFlake() {
	// 创建一个新的雪花算法生成器
	sf, err := NewSnowflake(1)
	if err != nil {
		panic(err)
	}
	fmt.Println("SnowFlake Init", sf)
	SF = sf
}

func GetSnowFlake() *Snowflake {
	return SF
}

// 创建新的 Snowflake 实例
func NewSnowflake(workerID int64) (*Snowflake, error) {
	epoch := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	workerIDBits := int64(10) // 机器 ID 占用 10 位
	sequenceBits := int64(12) // 序列号占用 12 位

	maxWorkerID := int64(-1 ^ (-1 << workerIDBits)) // 最大机器 ID (1023)
	maxSequence := int64(-1 ^ (-1 << sequenceBits)) // 最大序列号 (4095)

	if workerID < 0 || workerID > int64(maxWorkerID) {
		return nil, errors.New("workerID 超出范围")
	}

	return &Snowflake{
		timestamp:    0,
		workerID:     workerID,
		sequence:     0,
		epoch:        epoch,
		workerIDBits: workerIDBits,
		sequenceBits: sequenceBits,
		maxWorkerID:  maxWorkerID,
		maxSequence:  maxSequence,
		timeShift:    workerIDBits + sequenceBits,
		workerShift:  sequenceBits,
	}, nil
}

// 生成下一个唯一 ID
func (sf *Snowflake) NextID() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	now := time.Now().UnixMilli() // 当前时间戳（毫秒）
	if now < sf.timestamp {
		panic("时钟回拨，无法生成 ID")
	}

	if now == sf.timestamp {
		// 同一毫秒内递增序列号
		sf.sequence = (sf.sequence + 1) & sf.maxSequence
		if sf.sequence == 0 {
			// 如果序列号用尽，等待下一毫秒
			for now <= sf.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 如果是下一毫秒，重置序列号
		sf.sequence = 0
	}

	sf.timestamp = now

	// 生成 ID
	id := ((now - sf.epoch) << sf.timeShift) | (sf.workerID << sf.workerShift) | sf.sequence
	return id
}
