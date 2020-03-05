package domadoma

import "time"

type DealBase interface{
	CreateDeal()

	GetDeal()

	ListDeal()

	UpdateDeal()

	DeleteDeal()
}

type Deal struct {
	Id int `json:"id"`
	Food string `json:"name"`
	Quantity int `json:"quantity"`
	ConsumerId int `json:"consumer_id"`
	ProducerId int `json:"producer_id"`
	Created time.Time `json:"created"`
}
