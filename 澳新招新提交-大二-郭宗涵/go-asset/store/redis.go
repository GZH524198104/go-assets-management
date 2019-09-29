package store

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var redisPool *redis.Pool

func init() {
	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdle:         5,
		MaxConnLifetime: time.Second * 10,
		MaxActive:       100,
	}
}

func GetRedisConn() redis.Conn {
	return redisPool.Get()
}
