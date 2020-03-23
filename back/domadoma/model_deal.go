package domadoma

import "time"

type DealBase interface {
	CreateDeal(deal *Deal) (*Deal, error)

	GetDeal(id int) (*Deal, error)

	ListDeals() ([]*Deal, error)

	ListConsumerDeals(id int) ([]*Deal, error)

	ListProducerDeals(id int) ([]*Deal, error)

	UpdateDeal(id int, deal *Deal) (*Deal, error)

	DeleteDeal(id int) error
}


type Deal struct {
	Id int `json:"id" binding:"required"`
	Food string `json:"name" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
	ConsumerId int `json:"consumer_id" binding:"required"`
	ProducerId int `json:"producer_id" binding:"required"`
	Created time.Time `json:"created" binding:"required"`
	Finished time.Time `json:"finished,omitempty"`
	Active bool `json:"active"`
}
