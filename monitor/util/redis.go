package util

import (
	"time"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func NewRedisPool() *redis.Pool {
	hostPort := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)

	instance := &redis.Pool{
		MaxActive: 65536,
		MaxIdle: 256,
		IdleTimeout: 128 * time.Second,
		Wait: true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", hostPort)
			if err != nil {
				panic(fmt.Sprintf("fail to connect the redis server %s error %v\n", hostPort, err))
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return instance
}
