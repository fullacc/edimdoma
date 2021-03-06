package domadoma

type ConfigFile struct {
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
	RMQLogin     string `json:"rmq_login"`
	RMQPassword  string `json:"rmq_password"`
	RMQPort      string `json:"rmq_port"`
}

const (
	_ = iota //0
	BLunch //1
	Pervoe //2
	Vtoroe //3
	Salad //4
	Vypechka //5
	Desert //6
	Drugoe //7
)

const (
	_ = iota //0
	BezMyasa //1
	Govyadina //2
	Baranina //3
	Svinina //4
	Kurica //5
	Ptica //6
	Ryba //7
	More //8
	Baska //9
)

const(
	Null = iota //0
	Tru //1
	Fals //2
)

const(
	_ = iota //0
	NeOstro //1
	Slegka //2
	Ostrovato //3
	Podzhog //4
)
