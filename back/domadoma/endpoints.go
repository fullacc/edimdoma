package domadoma

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"net/http"
)

func renderError(w http.ResponseWriter,msg string,statuscode int) {
	w.WriteHeader(statuscode)
	w.Write([]byte(msg))
}

func getToken(token string) (*UserInfo,error) {
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