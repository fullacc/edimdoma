package OfferLog

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreOfferLogBase(db *pg.DB) (OfferLogBase, error) {
	//create schema
	for _, model := range []interface{}{(*OfferLog)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return nil, err
		}
	}
	return &postgreOfferLogBase{db: db}, nil
}

type postgreOfferLogBase struct {
	db *pg.DB
}

func (p *postgreOfferLogBase) CreateOfferLog(offerLog *OfferLog) (*OfferLog, error) {
	err := p.db.Insert(offerLog)
	if err != nil {
		return nil, err
	}
	return offerLog, nil
}

func (p *postgreOfferLogBase) GetOfferLog(offerLog *OfferLog) (*OfferLog, error) {
	err := p.db.Select(offerLog)
	if err != nil {
		return nil, err
	}
	return offerLog, nil
}

func (p *postgreOfferLogBase) ListOfferLogs() ([]*OfferLog, error) {
	var offerLogs []*OfferLog
	err := p.db.Model(&offerLogs).Select()
	if err != nil {
		return nil, err
	}
	return offerLogs, nil
}

func (p *postgreOfferLogBase) ListProducerOfferLogs(id int) ([]*OfferLog, error) {
	var offerLogs []*OfferLog
	err := p.db.Model(&offerLogs).Where("producer_id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return offerLogs, nil
}

func (p *postgreOfferLogBase) UpdateOfferLog(offerLog *OfferLog) (*OfferLog, error) {
	err := p.db.Update(offerLog)
	if err != nil {
		return nil, err
	}
	return offerLog, nil
}

func (p *postgreOfferLogBase) DeleteOfferLog(id int) error {
	offerLog := &OfferLog{Id: id}
	err := p.db.Delete(offerLog)
	if err != nil {
		return err
	}
	return nil
}
