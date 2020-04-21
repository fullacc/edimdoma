package Connection

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
)

func ConnectToDB(configfile *domadoma.ConfigFile) pg.DB{
	db := pg.Connect(&pg.Options{
		Database: configfile.PgDbName,
		Addr:     configfile.PgDbHost + ":" + configfile.PgDbPort,
		User:     configfile.PgDbUser,
		Password: configfile.PgDbPassword,
		})
	return *db
}