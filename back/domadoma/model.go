package domadoma

type ConfigFile struct{
	Port string `json:"port"`
	Host string `json:"host"`
	DbHost string `json:"dbhost"`
	DbPort string `json:"dbport"`
	Password string `json:"password"`
	User string `json:"user"`
	Name string `json:"name"`
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
	Token string
	Permission int
	UserId int
}

type UserRegister struct {
	UserName string `json:"user_name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Password string `json:"password"`
}

type UserLogin struct {
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}