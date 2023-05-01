package redis

import (
	"github.com/Zhoangp/Cart-Service/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cf *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cf.Redis.Address,
		Password: cf.Redis.Password, // no password set
		DB:       cf.Redis.Db,       // use default DB
	})
	return client, nil
}
