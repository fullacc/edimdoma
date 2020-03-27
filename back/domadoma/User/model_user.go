package User

type UserBase interface{
	CreateUser(user *User) (*User, error)

	GetUser(user *User) (*User, error)

	ListUsers() ([]*User, error)

	UpdateUser(user *User) (*User, error)

	DeleteUser(id int)  error
}

type User struct {
	Id int `json:"id,omitempty"`
	UserName string `json:"user_name" binding:"required"`
	PasswordHash string `json:"-"`
	Name string `json:"name" binding:"required"`
	Surname string `json:"surname" binding:"required"`
	RatingTotal float64 `json:"-"`
	RatingN float64 `json:"-"`
	Rating float64 `json:"rating"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required"`
	City string `json:"city" binding:"required"`
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
)

type UserInfo struct {
	Token      string
	Permission int
	UserId     int
}

type UserRegister struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserLogin struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserChangePassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}