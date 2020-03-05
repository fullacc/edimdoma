package domadoma

type UserBase interface{
	CreateUser()

	GetUser()

	ListUsers()

	UpdateUser()

	DeleteUser()
}

type User struct {
	Id int `json:"id"`
	UserName string `json:"user_name"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Phone string `json:"phone"`
	City string `json:"city"`
}