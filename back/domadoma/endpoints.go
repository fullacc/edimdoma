package domadoma

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

func GetToken(token string) (*UserInfo,error) {
	if len(token) == 0 {
		return nil,errors.New("no token provided")
	}
	redisClient := Connect()
	val, err := redisClient.Get(token).Result()
	uInfo := &UserInfo{}
	if err == redis.Nil {
		return nil,errors.New("No such token")
	} else if err != nil {
		return nil,err
	} else {
		err := json.Unmarshal([]byte(val), uInfo)
		if err != nil {
			return nil,err
		}
	}
	return uInfo,nil
}

func Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}