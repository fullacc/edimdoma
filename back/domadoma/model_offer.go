package domadoma

import (
	"time"
)

type OfferBase interface{
	CreateOffer()

	GetOffer()

	ListOffer()

	UpdateOffer()

	DeleteOffer()
}

type Offer struct {
	Id int `json:"id"`
	ProducerId int `json:"producer_id"`
	Food string `json:"name"`
	Price int `json:"price"`
	InitialQuantity int `json:"initial_quantity"`
	AvailableQuantity int `json:"available_quantity"`
	Location []float64 `json:"location"`
	Created time.Time `json:"created"`
}
