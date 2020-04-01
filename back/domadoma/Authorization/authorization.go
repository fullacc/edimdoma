package Authorization

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"
)

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

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateToken() (string, error) {
	b, err := GenerateRandomBytes(32)
	return base64.URLEncoding.EncodeToString(b), err
}
