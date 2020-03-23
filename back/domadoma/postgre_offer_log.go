package domadoma

import "github.com/go-pg/pg"

func NewPostgreOfferLogBase(configfile *ConfigFile) (OfferLogBase, error) {

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


func (p *postgreBase) CreateOfferLog(offerLog *OfferLog) (*OfferLog, error) {
	err := p.db.Insert(offerLog)
	if err != nil {
		return nil,err
	}
	return offerLog,nil
}

func (p *postgreBase) GetOfferLog(id int) (*OfferLog, error) {
	offerLog := &OfferLog{Id:id}
	err := p.db.Select(&offerLog)
	if err != nil {
		return nil, err
	}
	return offerLog, nil
}

func (p *postgreBase) ListOfferLogs() ([]*OfferLog, error) {
	var offerLogs []*OfferLog
	err := p.db.Select(offerLogs)
	if err != nil {
		return nil, err
	}
	return offerLogs,nil
}

func (p *postgreBase) UpdateOfferLog(id int, offerLog *OfferLog) (*OfferLog, error) {
	offerLog1 := &OfferLog{Id:id}
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

func (p *postgreBase) DeleteOfferLog(id int) error {
	offerLog := &OfferLog{Id: id}
	err := p.db.Delete(offerLog)
	if err != nil {
		return err
	}
	return nil
}