package domadoma

import (
	"github.com/fullacc/edimdoma/back/domadoma/Deal"
	"github.com/fullacc/edimdoma/back/domadoma/Feedback"
	"github.com/fullacc/edimdoma/back/domadoma/Offer"
	"github.com/fullacc/edimdoma/back/domadoma/Request"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type postgreBase struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Deal.Deal)(nil),(*Feedback.Feedback)(nil),(*Offer.Offer)(nil),(*OffersLog)(nil),(*Request.Request)(nil),(*User.User)(nil)} {
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

