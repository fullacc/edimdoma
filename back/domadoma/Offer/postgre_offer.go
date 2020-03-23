package Offer

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
)

func NewPostgreOfferBase(configfile *domadoma.ConfigFile) (OfferBase, error) {

	db := pg.Connect(&pg.Options{
		Database: configfile.Name,
		Addr: configfile.DbHost + ":" + configfile.DbPort,
		User: "postgres",
		Password: configfile.Password,
	})

	err := domadoma.createSchema(db)
	if err != nil {
		return nil, err
	}
	return &domadoma.postgreBase{db: db}, nil
}


func (p *domadoma.postgreBase) CreateOffer(offer *Offer) (*Offer, error) {
	err := p.db.Insert(offer)
	if err != nil {
		return nil,err
	}
	return offer, nil
}

func (p *domadoma.postgreBase) GetOffer(id int) (*Offer, error) {
	offer := &Offer{Id: id}
	err := p.db.Select(&offer)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (p *domadoma.postgreBase) ListOffers() ([]*Offer, error) {
	var offers []*Offer
	err := p.db.Select(offers)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (p *domadoma.postgreBase) ListProducerOffers(id int) ([]*Offer, error){
	var offers []*Offer
	err := p.db.Model(&offers).Where("Producer_Id = ?",id).Select()
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (p *domadoma.postgreBase) UpdateOffer(id int, offer *Offer) (*Offer, error) {
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

func (p *domadoma.postgreBase) DeleteOffer(id int) error {
	offer := &Offer{Id: id}
	err := p.db.Delete(offer)
	if err != nil {
		return err
	}
	return nil
}