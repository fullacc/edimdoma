package domadoma

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/go-redis/redis"
	"regexp"
)

func GetToken(token string) (*User.UserInfo,error) {
	if len(token) == 0 {
		return nil,errors.New("no token provided")
	}
	redisClient := Connect()
	val, err := redisClient.Get(token).Result()
	uInfo := &User.UserInfo{}
	if err == redis.Nil {
		return nil,errors.New("no such token")
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

func Validator(vid int,s string) (bool,error) {
	usernameregex := `^([\.\_a-z0-9]{1,30})$`
	phoneregex := `^[0-9]{10}$`
	emailregex := `^([a-z0-9_\-\.]+)@([a-z0-9_\-\.]+)\.([a-z]{2,5})$`
	switch vid{
	case User.Usrnm:
		matched, err := regexp.Match(usernameregex, []byte(s))
		if err != nil {
			return false,err
		}
		if !matched {
			return false, nil
		}
		matched, err = regexp.Match(`([a-z]{1,30})`, []byte(s))
		if err != nil {
			return false,err
		}
		if !matched {
			return false, nil
		}
		if s[0]=='.' || s[len(s)-1]=='.'{
			return false, nil
		}
		for i,v := range s{
			if v=='.' && s[i+1]=='.'{
				return false,nil
			}
		}
		return true,nil
	case User.Phn:
		matched, err := regexp.Match(phoneregex, []byte(s))
		if err != nil {
			return false,err
		}
		if !matched {
			return false, nil
		}
		return true,nil
	case User.Eml:
		matched, err := regexp.Match(emailregex, []byte(s))
		if err != nil {
			return false,err
		}
		if !matched {
			return false, nil
		}
		return true,nil
	default:
		return false,errors.New("unknown tag")
	}
}
