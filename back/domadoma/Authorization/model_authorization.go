package Authorization

import (
	"errors"
	"regexp"
	"time"
)

type AuthorizationBase interface{
	GetAuthToken(token string) (*AuthToken, error)

	GetRegistrationToken(token string) (*RegistrationToken, error)

	SetToken(token string, data []byte, expr time.Duration) error

	DeleteToken(token string) error
}

const (
	Unknown = iota
	Admin
	Manager
	Regular
)

const (
	Usrnm = iota
	Phn
	Eml
	Cd
)

type AuthToken struct {
	Token      string
	Permission int
	UserId     int
}

type RegistrationToken struct {
	Token string `json:"token"`
	Phone string `json:"phone" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type RegistrationPhone struct {
	Phone string `json:"phone"`
}

type RegistrationCode struct {
	Code string `json:"code"`
}

type UserRegister struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserChangePassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func Validator(vid int, s string) (bool, error) {
	usernameregex := `^([\.\_a-z0-9]{1,30})$`
	phoneregex := `^[0-9]{10}$`
	coderegex := `^[0-9]{6}$`
	emailregex := `^([a-z0-9_\-\.]+)@([a-z0-9_\-\.]+)\.([a-z]{2,5})$`
	switch vid {
	case Usrnm:
		matched, err := regexp.Match(usernameregex, []byte(s))
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
		matched, err = regexp.Match(`([a-z]{1,30})`, []byte(s))
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
		if s[0] == '.' || s[len(s)-1] == '.' {
			return false, nil
		}
		for i, v := range s {
			if v == '.' && s[i+1] == '.' {
				return false, nil
			}
		}
		return true, nil
	case Phn:
		matched, err := regexp.Match(phoneregex, []byte(s))
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
		return true, nil
	case Eml:
		matched, err := regexp.Match(emailregex, []byte(s))
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
		return true, nil
	case Cd:
		matched, err := regexp.Match(coderegex, []byte(s))
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
		return true, nil
	default:
		return false, errors.New("unknown tag")
	}
}
