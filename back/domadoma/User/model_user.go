package User

type UserBase interface {
	CreateUser(user *User) (*User, error)

	GetUser(user *User) (*User, error)

	ListUsers() ([]*User, error)

	UpdateUser(user *User) (*User, error)

	DeleteUser(id int) error
}

type User struct {
	Id           int       `json:"id,omitempty"`
	UserName     string    `json:"user_name" binding:"required"`
	Role         int       `json:"-"`
	PasswordHash []byte    `json:"-"`
	Name         string    `json:"name" binding:"required"`
	Surname      string    `json:"surname" binding:"required"`
	RatingTotal  float64   `json:"-"`
	RatingN      float64   `json:"-"`
	Rating       float64   `json:"rating"`
	Phone        string    `json:"phone" binding:"required"`
	City         string    `json:"city" binding:"required"`
	Location     []float64 `json:"location"`
}
