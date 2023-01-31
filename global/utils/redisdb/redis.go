package redisdb

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisInterface interface {
	Client() *redis.Client
}

type redisStruct struct {
	client *redis.Client
}

func InitRedis(host string, port string, password string, database string) RedisInterface {
	dbInt, _ := strconv.Atoi(database)
	serverString := fmt.Sprintf("%s:%s", host, port)
	var rdb *redis.Client

	if len(password) > 0 {
		rdb = redis.NewClient(&redis.Options{
			Addr:     serverString,
			DB:       dbInt,
			Password: password,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr: serverString,
			DB:   dbInt,
		})
	}

	return &redisStruct{
		client: rdb,
	}
}

func (r *redisStruct) Client() *redis.Client {
	return r.client
}
