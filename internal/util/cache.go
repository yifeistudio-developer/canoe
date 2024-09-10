package util

import (
	"context"
	"github.com/bluele/gcache"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	rdb   *redis.Client
	cache gcache.Cache
	ctx   = context.Background()
)

func newCache() {
	if rdb == nil {
		rdb = newRedisClient()
	}
	cache = gcache.New(20).
		LRU().
		Expiration(5 * time.Minute).
		LoaderFunc(func(i interface{}) (interface{}, error) {
			if rdb == nil {
				rdb = newRedisClient()
			}
			return rdb.Get(ctx, i.(string)).Result()
		}).Build()
}

func newRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 10, // 最大连接数
	})
	return rdb
}
