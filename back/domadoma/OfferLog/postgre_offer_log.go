package OfferLog

import (
	"../../domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreOfferLogBase(configfile *domadoma.ConfigFile) (OfferLogBase, error) {

	db := pg.Connect(&pg.Options{
		Database: configfile.PgDbName,
		Addr: configfile.PgDbHost + ":" + configfile.PgDbPort,
		User: configfile.PgDbUser,
		Password: configfile.PgDbPassword,
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

func (p *postgreOfferLogBase) ListProducerOfferLogs(id int) ([]*OfferLog, error){
	var offerLogs []*OfferLog
	err := p.db.Model(&offerLogs).Where("Producer_Id = ?",id).Select()
	if err != nil {
		return nil, err
	}
	return offerLogs, nil
}

func (p *postgreOfferLogBase) UpdateOfferLog(offerLog *OfferLog) (*OfferLog, error) {
	err := p.db.Update(offerLog)
	if err != nil {
		return nil,err
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