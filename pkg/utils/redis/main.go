package redis

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}
