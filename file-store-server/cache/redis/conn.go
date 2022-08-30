package redis

import (
	"time"

	rs "github.com/go-redis/redis/v8"
)

var (
	redisCli *rs.Client
)

func newRedisPool() *rs.Client {
	return rs.NewClient(&rs.Options{
		Network:            "tcp",
		Addr:               "127.0.0.1:6379",
		PoolSize:           15,
		MinIdleConns:       10,
		IdleTimeout:        300 * time.Second,
		IdleCheckFrequency: 60 * time.Second,
	})
}

func init() {
	redisCli = newRedisPool()
}

func RedisPool() *rs.Client {
	return redisCli
}
