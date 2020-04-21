package Offer

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreOfferBase(db *pg.DB) (OfferBase, error) {
	//create schema
	for _, model := range []interface{}{(*Offer)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return nil, err
		}
	}
	return &postgreOfferBase{db: db}, nil
}

type postgreOfferBase struct {
	db *pg.DB
}

func (p *postgreOfferBase) CreateOffer(offer *Offer) (*Offer, error) {
	err := p.db.Insert(offer)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (p *postgreOfferBase) GetOffer(offer *Offer) (*Offer, error) {
	err := p.db.Select(offer)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (p *postgreOfferBase) ListOffers() ([]*Offer, error) {
	var offers []*Offer
	err := p.db.Model(&offers).Select()
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (p *postgreOfferBase) ListProducerOffers(id int) ([]*Offer, error) {
	var offers []*Offer
	err := p.db.Model(&offers).Where("producer_id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (p *postgreOfferBase) UpdateOffer(offer *Offer) (*Offer, error) {
	err := p.db.Update(offer)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (p *postgreOfferBase) DeleteOffer(id int) error {
	offer := &Offer{Id: id}
	err := p.db.Delete(offer)
	if err != nil {
		return err
	}
	return nil
}
