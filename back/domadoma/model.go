package domadoma

type ConfigFile struct{
	ApiPort      string `json:"apiport"`
	ApiHost      string `json:"apihost"`
	PgDbHost     string `json:"pgdbhost"`
	PgDbPort     string `json:"pgdbport"`
	PgDbPassword string `json:"pgdbpassword"`
	PgDbUser     string `json:"pgdbuser"`
	PgDbName     string `json:"pgdbname"`
	RdHost       string `json:"rdhost"`
	RdPort       string `json:"rdport"`
	RdPass       string `json:"rdpass"`
	SMSlogin     string `json:"smslogin"`
	SMSpass      string `json:"smspass"`
}

