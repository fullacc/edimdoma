package Offer

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreOfferBase(configfile *domadoma.ConfigFile) (OfferBase, error) {

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
	return &postgreOfferBase{db: db}, nil
}

type postgreOfferBase struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Offer)(nil)} {
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

func (p *postgreOfferBase) CreateOffer(offer *Offer) (*Offer, error) {
	err := p.db.Insert(offer)
	if err != nil {
		return nil,err
	}
	return offer, nil
}

func (p *postgreOfferBase) GetOffer(id int) (*Offer, error) {
	offer := &Offer{Id: id}
	err := p.db.Select(&offer)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (p *postgreOfferBase) ListOffers() ([]*Offer, error) {
	var offers []*Offer
	err := p.db.Select(offers)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (p *postgreOfferBase) ListProducerOffers(id int) ([]*Offer, error){
	var offers []*Offer
	err := p.db.Model(&offers).Where("Producer_Id = ?",id).Select()
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (p *postgreOfferBase) UpdateOffer(id int, offer *Offer) (*Offer, error) {
	offer1 := &Offer{Id: id}
	err := p.db.Select(offer1)
	if err != nil {
		return nil,err
	}
	offer1 = offer
	err = p.db.Update(offer1)
	if err != nil {
		return nil,err
	}
	return offer1, nil
}

func (p *postgreOfferBase) DeleteOffer(id int) error {
	offer := &Offer{Id: id}
	err := p.db.Delete(offer)
	if err != nil {
		return err
	}
	return nil
}