package Authorization

import (
	"github.com/segmentio/encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

func NewRedisAuthorizationBase(client *redis.Client) (AuthorizationBase, error) {
	_, err := client.Ping().Result()
	if err !=nil {
		return nil, err
	}
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

func (r redisAuthorizationBase) GetRegistrationToken(token string) (*RegistrationToken, error) {
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

func (r redisAuthorizationBase) GetForgotToken(token string) (*ForgotToken, error) {
	if len(token) == 0 {
		return nil, errors.New("no token provided")
	}
	val, err := r.rdb.Get(token).Result()
	uInfo := &ForgotToken{}
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

func (r redisAuthorizationBase) SetToken(token string, data []byte, expr time.Duration) error {
	return r.rdb.Set(token, data, expr).Err()
}

func (r redisAuthorizationBase) DeleteToken(token string) error {
	return r.rdb.Del(token).Err()
}
