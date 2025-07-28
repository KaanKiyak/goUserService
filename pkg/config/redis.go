package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

var Rdb *redis.Client
var Ctx = context.Background()

func RedisConnect() {

	reddisAddr := os.Getenv("REDIS_ADDR")
	reddisDbStr := os.Getenv("REDIS_DB")

	reddisDb, err := strconv.Atoi(reddisDbStr)
	if err != nil {
		log.Fatal(err)
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr: reddisAddr,
		DB:   reddisDb,
	})
	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatal(err)
	}
	log.Println("Redis connected", reddisAddr)

}
