package init

import (
	"fmt"
	"log"

	"github.com/firmfoundation/dbquery/config"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ConnectRedis(config *config.DbConfig) bool {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Println("Failed to connect to the Redis", err)
		return false
	}

	fmt.Println("? Connected Successfully to the Redis")
	return true
}
