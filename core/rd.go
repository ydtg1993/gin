package core

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var Redis *redis.Client

func init() {
	// Initialize Redis client
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Config.GetString("redis.host"), Config.GetInt("redis.port")),
		Password: Config.GetString("redis.password"),
		DB:       Config.GetInt("redis.db"),
		PoolSize: Config.GetInt("redis.pool_size"),
	})

	// Ping the Redis server to check if the connection is successful
	pong, err := Redis.Ping(Redis.Context()).Result()
	if err != nil {
		panic(fmt.Sprintf("unable to connect to Redis: %s", err))
	}
	fmt.Printf("Connected to Redis, Ping Response: %s\n", pong)

	// Set a default connection timeout
	Redis.Options().ReadTimeout = Config.GetDuration("redis.timeout") * time.Second
	Redis.Options().WriteTimeout = Config.GetDuration("redis.timeout") * time.Second
}
