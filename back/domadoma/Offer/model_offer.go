package Offer

import (
	"time"
)

type OfferBase interface {
	CreateOffer(offer *Offer) (*Offer, error)

	GetOffer(offer *Offer) (*Offer, error)

	ListOffers() ([]*Offer, error)

	ListProducerOffers(id int) ([]*Offer, error)

	UpdateOffer(offer *Offer) (*Offer, error)

	DeleteOffer(id int) error
}

type Offer struct {
	Id                int       `json:"id"`
	ProducerId        int       `json:"producer_id"`
	FoodName          string    `json:"food_name" binding:"required"`
	Type              int       `json:"type"`
	Myaso             int       `json:"myaso"`
	Halal             int       `json:"halal"`
	Vegan             int       `json:"vegan"`
	Spicy             int       `json:"spicy"`
	Description       string    `json:"description"`
	Price             int       `json:"price" binding:"required"`
	InitialQuantity   int       `json:"initial_quantity" binding:"required"`
	AvailableQuantity int       `json:"available_quantity"`
	Location          []float64 `json:"location" binding:"required"`
	Created           time.Time `json:"created"`
}
