package domadoma

import "net/http"

type UserEndpoints interface{
	GetUser() func(w http.ResponseWriter, r *http.Request)

	CreateUser() func(w http.ResponseWriter, r *http.Request)

	ListUsers() func(w http.ResponseWriter, r *http.Request)

	UpdateUser() func(w http.ResponseWriter, r *http.Request)

	DeleteUser() func(w http.ResponseWriter, r *http.Request)

}

func NewUserEndpoints(userBase UserBase) UserEndpoints {
	return &UserEndpointsFactory{userBase: userBase}
}

type UserEndpointsFactory struct{
	userBase UserBase
}