package dao

import (
	"context"

	config "gamesvr/config"

	"github.com/redis/go-redis/v9"
)

type RedisMgr struct {
	RedisClient *redis.Client
	ctx         context.Context
}

var RedisMgrInstance *RedisMgr

func init() {
	RedisMgrInstance = &RedisMgr{}
	RedisMgrInstance.ctx = context.Background()

	RedisMgrInstance.Init()
}

func (r *RedisMgr) Init() {
	r.RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})

	_, err := r.RedisClient.Ping(r.ctx).Result()
	if err != nil {
		panic(err)
	}
}

func (r *RedisMgr) Close() {
	if r.RedisClient != nil {
		r.RedisClient.Close()
	}
}

func (r *RedisMgr) GetRedisClient() *redis.Client {
	return r.RedisClient
}

func (r *RedisMgr) Set(key string, value interface{}) error {
	return r.RedisClient.Set(r.ctx, key, value, 0).Err()
}

func (r *RedisMgr) Get(key string) (string, error) {
	return r.RedisClient.Get(r.ctx, key).Result()
}
func (r *RedisMgr) Del(key string) error {
	return r.RedisClient.Del(r.ctx, key).Err()
}
func (r *RedisMgr) HSet(key, field string, value interface{}) error {
	return r.RedisClient.HSet(r.ctx, key, field, value).Err()
}

func (r *RedisMgr) HGet(key, field string) (string, error) {
	return r.RedisClient.HGet(r.ctx, key, field).Result()
}

func (r *RedisMgr) HGetAll(key string) (map[string]string, error) {
	return r.RedisClient.HGetAll(r.ctx, key).Result()
}

func (r *RedisMgr) HDel(key string, fields ...string) error {
	return r.RedisClient.HDel(r.ctx, key, fields...).Err()
}

// HExists reports whether the specified field exists in the hash stored at key. if the key does not exist, it returns false.
func (r *RedisMgr) HExists(key, field string) (bool, error) {
	return r.RedisClient.HExists(r.ctx, key, field).Result()
}
