package Deal

import (
	"fmt"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreDealBase(configfile *domadoma.ConfigFile) (DealBase, error) {
	fmt.Println(configfile)
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
	return &postgreDealBase{db: db}, nil
}

type postgreDealBase struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Deal)(nil)} {
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

func (p *postgreDealBase) CreateDeal(deal *Deal,) (*Deal, error) {
	err := p.db.Insert(deal)
	if err != nil {
		return nil,err
	}
	return deal,nil
}

func (p *postgreDealBase) GetDeal(deal *Deal) (*Deal, error) {
	err := p.db.Select(deal)
	if err != nil {
		return nil, err
	}
	return deal, nil
}

func (p *postgreDealBase) ListDeals() ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Model(&deals).Select()
	if err != nil {
		return nil, err
	}
	return deals,nil
}

func (p *postgreDealBase) ListConsumerDeals(id int) ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Model(&deals).Where("consumer_id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return deals, nil
}

func (p *postgreDealBase) ListProducerDeals(id int) ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Model(&deals).Where("producer_id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return deals, nil
}
func (p *postgreDealBase) ListActiveDeals() ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Model(&deals).Where("complete = ?", "false").Select()
	if err != nil {
		return nil, err
	}
	return deals,nil
}

func (p *postgreDealBase) ListActiveConsumerDeals(id int) ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Model(&deals).Where("consumer_Id = ?", id).Where("complete = ?", "false").Select()
	if err != nil {
		return nil, err
	}
	return deals, nil
}

func (p *postgreDealBase) ListActiveProducerDeals(id int) ([]*Deal, error) {
	var deals []*Deal
	err := p.db.Model(&deals).Where("producer_Id = ?", id).Where("complete = ?", "false").Select()
	if err != nil {
		return nil, err
	}
	return deals, nil
}

func (p *postgreDealBase) UpdateDeal(deal *Deal) (*Deal, error) {
	err := p.db.Update(deal)
	if err != nil {
		return nil,err
	}
	return deal, nil
}

func (p *postgreDealBase) DeleteDeal(id int) error {
	deal := &Deal{Id: id}
	err := p.db.Delete(deal)
	if err != nil {
		return err
	}
	return nil
}