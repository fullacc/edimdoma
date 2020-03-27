package Deal

import "time"

type DealBase interface {
	CreateDeal(deal *Deal) (*Deal, error)

	GetDeal(id int) (*Deal, error)

	ListDeals() ([]*Deal, error)

	ListConsumerDeals(id int) ([]*Deal, error)

	ListProducerDeals(id int) ([]*Deal, error)

	ListActiveDeals() ([]*Deal, error)

	ListActiveConsumerDeals(id int) ([]*Deal, error)

	ListActiveProducerDeals(id int) ([]*Deal, error)

	UpdateDeal(deal *Deal) (*Deal, error)

	DeleteDeal(id int) error
}


type Deal struct {
	Id int `json:"id"`
	Food string `json:"name"`
	Price int `json:"price" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
	ConsumerId int `json:"consumer_id"`
	ProducerId int `json:"producer_id"`
	Created time.Time `json:"created"`
	Finished time.Time `json:"finished,omitempty"`
	Complete bool `json:"complete"`
}
