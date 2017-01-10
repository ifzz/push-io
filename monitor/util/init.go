package util

import (
    "sync"
    "github.com/garyburd/redigo/redis"
)

var config *Config
var once sync.Once
var redisPool *redis.Pool

func init() {
    once.Do(func() {
        config = NewConfig()
        redisPool = NewRedisPool()
    })
}
