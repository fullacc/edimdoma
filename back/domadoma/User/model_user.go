package User

type UserBase interface{
	CreateUser(user *User) (*User, error)

	GetUser(user *User) (*User, error)

	ListUsers() ([]*User, error)

	UpdateUser(id int, user *User) (*User, error)

	DeleteUser(id int)  error
}

type User struct {
	Id int `json:"id,omitempty"`
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"-"`
	PasswordHash string
	Name string `json:"name" binding:"required"`
	Surname string `json:"surname" binding:"required"`
	RatingTotal float64
	RatingN float64
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required"`
	City string `json:"city" binding:"required"`
}