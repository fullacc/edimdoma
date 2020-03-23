package domadoma

import "github.com/go-pg/pg"

func NewPostgreOfferBase(configfile *ConfigFile) (OfferBase, error) {

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


func (p *postgreBase) CreateOffer(offer *Offer) (*Offer, error) {
	err := p.db.Insert(offer)
	if err != nil {
		return nil,err
	}
	return offer, nil
}

func (p *postgreBase) GetOffer(id int) (*Offer, error) {
	offer := &Offer{Id:id}
	err := p.db.Select(&offer)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (p *postgreBase) ListOffers() ([]*Offer, error) {
	var offers []*Offer
	err := p.db.Select(offers)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (p *postgreBase) ListProducerOffers(id int) ([]*Offer, error){
	var offers []*Offer
	err := p.db.Model(&offers).Where("Producer_Id = ?",id).Select()
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (p *postgreBase) UpdateOffer(id int, offer *Offer) (*Offer, error) {
	offer1 := &Offer{Id:id}
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

func (p *postgreBase) DeleteOffer(id int) error {
	offer := &Offer{Id: id}
	err := p.db.Delete(offer)
	if err != nil {
		return err
	}
	return nil
}