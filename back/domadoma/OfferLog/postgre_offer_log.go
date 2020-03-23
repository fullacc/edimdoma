package OfferLog

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreOfferLogBase(configfile *domadoma.ConfigFile) (OfferLogBase, error) {

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
	return &postgreOfferLogBase{db: db}, nil
}

type postgreOfferLogBase struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*OfferLog)(nil)} {
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

func (p *postgreOfferLogBase) CreateOfferLog(offerLog *OfferLog) (*OfferLog, error) {
	err := p.db.Insert(offerLog)
	if err != nil {
		return nil,err
	}
	return offerLog,nil
}

func (p *postgreOfferLogBase) GetOfferLog(id int) (*OfferLog, error) {
	offerLog := &OfferLog{Id: id}
	err := p.db.Select(&offerLog)
	if err != nil {
		return nil, err
	}
	return offerLog, nil
}

func (p *postgreOfferLogBase) ListOfferLogs() ([]*OfferLog, error) {
	var offerLogs []*OfferLog
	err := p.db.Select(offerLogs)
	if err != nil {
		return nil, err
	}
	return offerLogs,nil
}

func (p *postgreOfferLogBase) UpdateOfferLog(id int, offerLog *OfferLog) (*OfferLog, error) {
	offerLog1 := &OfferLog{Id: id}
	err := p.db.Select(offerLog1)
	if err != nil {
		return nil,err
	}
	offerLog1 = offerLog
	err = p.db.Update(offerLog1)
	if err != nil {
		return nil,err
	}
	return offerLog1, nil
}

func (p *postgreOfferLogBase) DeleteOfferLog(id int) error {
	offerLog := &OfferLog{Id: id}
	err := p.db.Delete(offerLog)
	if err != nil {
		return err
	}
	return nil
}