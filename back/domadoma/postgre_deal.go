package domadoma

import "github.com/go-pg/pg"

func NewPostgreDealBase(configfile *ConfigFile) (DealBase, error) {

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


func (p *postgreBase) CreateDeal(deal *Deal,) (*Deal, error) {
	err := p.db.Insert(deal)
	if err != nil {
		return nil,err
	}
	return deal,nil
}

func (p *postgreBase) GetDeal(id int) (*Deal, error) {
	deal := &Deal{Id:id}
	err := p.db.Select(deal)
	if err != nil {
		return nil, err
	}
	return deal, nil
}

func (p *postgreBase) ListDeals() ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Select(deals)
	if err != nil {
		return nil, err
	}
	return deals,nil
}

func (p *postgreBase) ListConsumerDeals(id int) ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Model(&deals).Where("Consumer_Id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return deals, nil
}

func (p *postgreBase) ListProducerDeals(id int) ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Model(&deals).Where("Producer_Id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return deals, nil
}

func (p *postgreBase) UpdateDeal(id int, deal *Deal) (*Deal, error) {
	deal1 := &Deal{Id:id}
	err := p.db.Select(deal1)
	if err != nil {
		return nil,err
	}
	deal1 = deal
	err = p.db.Update(deal1)
	if err != nil {
		return nil,err
	}
	return deal1, nil
}

func (p *postgreBase) DeleteDeal(id int) error {
	deal := &Deal{Id: id}
	err := p.db.Delete(deal)
	if err != nil {
		return err
	}
	return nil
}