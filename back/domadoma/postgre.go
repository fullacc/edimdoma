package domadoma

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type postgreBase struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Consumer)(nil),(*Deal)(nil),(*Feedback)(nil),(*Offer)(nil),(*OffersLog)(nil),(*Producer)(nil),(*Request)(nil),(*User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

