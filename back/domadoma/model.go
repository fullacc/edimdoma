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
