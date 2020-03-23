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

type UserInfo struct {
	Token string
	Permission int
	UserId int
}

type UserLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}