package Authorization

import (
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

