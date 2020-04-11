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
}

const (
	Undef = iota //0
	BLunch //1
	Pervoe //2
	Vtoroe //3
	Salad //4
	Vypechka //5
	Desert //6
	Drugoe //7
)

const (
	Nolik = iota //0
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
	Hz = iota //0
	Slegka //1
	Ostrovato //2
	Podzhog //3
)
