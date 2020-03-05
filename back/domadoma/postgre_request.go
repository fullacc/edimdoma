package domadoma

import "github.com/go-pg/pg"

func NewPostgreRequestBase(configfile *ConfigFile) (RequestBase, error) {

	db := pg.Connect(&pg.Options{
		Database: configfile.Name,
		Addr: configfile.DbHost + ":" + configfile.DbPort,
		User: "postgres",
		Password: configfile.Password,
	})

	err := createSchema(db)
	if err != nil {
		return nil, err
	}
	return &postgreBase{db: db}, nil
}