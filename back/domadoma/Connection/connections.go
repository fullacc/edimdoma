package Connection

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
)

func ConnectToPostgre(configfile *domadoma.ConfigFile) *pg.DB{
	db := pg.Connect(&pg.Options{
		Database: configfile.PgDbName,
		Addr:     configfile.PgDbHost + ":" + configfile.PgDbPort,
		User:     configfile.PgDbUser,
		Password: configfile.PgDbPassword,
		})
	return db
}

func ConnectToRedis(configfile *domadoma.ConfigFile) *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr:     configfile.RdHost + ":" + configfile.RdPort,
		Password: configfile.RdPass, // no password set
		DB:       0,                 // use default DB
	})
	return client
}