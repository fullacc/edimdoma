package domadoma

import (
	"./Deal"
	"./Feedback"
	"./Offer"
	"./OfferLog"
	"./Request"
	"./User"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type PostgreBase struct {
	db *pg.DB
}

func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Deal.Deal)(nil),(*Feedback.Feedback)(nil),(*Offer.Offer)(nil),(*OfferLog.OfferLog)(nil),(*Request.Request)(nil),(*User.User)(nil)} {
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

