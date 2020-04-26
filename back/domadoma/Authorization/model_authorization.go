package Authorization

import (
	"time"
)

type AuthorizationBase interface {
	GetAuthToken(token string) (*AuthToken, error)

	GetRegistrationToken(token string) (*RegistrationToken, error)

	GetForgotToken(token string) (*ForgotToken, error)

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
	Token      string `json:"token"`
	Permission int `json:"permission" binding:"required"`
	UserId     int `json:"user_id" binding:"required"`
}

type RegistrationToken struct {
	Token string `json:"token"`
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type Phone struct {
	Phone string `json:"phone"`
}

type Code struct {
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

type ForgotLogin struct{
	Login string `json:"login" binding:"required"`
}

type ForgotToken struct{
	Token  string `json:"token"`
	UserId int    `json:"user_id"`
	Code   string `json:"code"`
}

type Password struct{
	Password string `json:"password" binding:"required"`
}