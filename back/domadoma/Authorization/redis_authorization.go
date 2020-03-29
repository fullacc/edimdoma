package Authorization

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-redis/redis"
	"time"
)

func NewRedisAuthorizationBase(configfile *domadoma.ConfigFile) (AuthorizationBase, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     configfile.RdHost+":"+configfile.RdPort,
		Password: configfile.RdPass, // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return &redisAuthorizationBase{rdb: client}, err
}

type redisAuthorizationBase struct {
	rdb *redis.Client
}


func (r redisAuthorizationBase) GetAuthToken(token string) (*AuthToken, error) {
	if len(token) == 0 {
		return nil, errors.New("no token provided")
	}
	val, err := r.rdb.Get(token).Result()
	uInfo := &AuthToken{}
	if err == redis.Nil {
		return nil, errors.New("no such token")
	} else if err != nil {
		return nil, err
	} else {
		err := json.Unmarshal([]byte(val), uInfo)
		if err != nil {
			return nil, err
		}
	}
	return uInfo, nil
}

func (r redisAuthorizationBase) GetRegisterToken(token string) (*RegistrationToken,error) {
	if len(token) == 0 {
		return nil, errors.New("no token provided")
	}
	val, err := r.rdb.Get(token).Result()
	uInfo := &RegistrationToken{}
	if err == redis.Nil {
		return nil, errors.New("no such token")
	} else if err != nil {
		return nil, err
	} else {
		err := json.Unmarshal([]byte(val), uInfo)
		if err != nil {
			return nil, err
		}
	}
	return uInfo, nil
}

func (r redisAuthorizationBase) SetToken(token string,data []byte, expr time.Duration) error {
	return r.rdb.Set(token,data,expr).Err()
}

func (r redisAuthorizationBase) DeleteToken(token string) error {
	return r.rdb.Del(token).Err()
}