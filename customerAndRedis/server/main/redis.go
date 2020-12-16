package main

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

//定义一个全局的redisPool

var pool *redis.Pool

//这个initPool函数要在main.go中的main函数中一开始就初始化好
//不会自动调用的
func initPool(address string, maxIdle int, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
