package auth

import (
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func Redis() *redis.Client {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	pwd := os.Getenv("rd_pwd")
	host := os.Getenv("rd_endpoint")
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd,
		DB:       0,
	})
	return rdb
}
