package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}

func redisKey(shortCode string) string {
	return "url:" + shortCode
}

func GetURL(shortCode string) (string, error) {
	return Rdb.Get(Ctx, redisKey(shortCode)).Result()
}

func SetURL(shortCode, longURL string) error {
	return Rdb.Set(
		Ctx,
		redisKey(shortCode),
		longURL,
		24*time.Hour,
	).Err()
}
