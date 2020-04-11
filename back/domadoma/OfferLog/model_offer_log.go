package OfferLog

import "time"

type OfferLogBase interface {
	CreateOfferLog(offerLog *OfferLog) (*OfferLog, error)

	GetOfferLog(offerLog *OfferLog) (*OfferLog, error)

	ListOfferLogs() ([]*OfferLog, error)

	ListProducerOfferLogs(id int) ([]*OfferLog, error)

	UpdateOfferLog(offerLog *OfferLog) (*OfferLog, error)

	DeleteOfferLog(id int) error
}

type OfferLog struct {
	Id                int       `json:"id"`
	ProducerId        int       `json:"producer_id" binding:"required"`
	FoodName          string    `json:"food_name" binding:"required"`
	Type              int       `json:"type"`
	Myaso             int       `json:"myaso"`
	Halal             int       `json:"halal"`
	Vegan             int       `json:"vegan"`
	Spicy             int       `json:"spicy"`
	Description       string    `json:"description"`
	Price             int       `json:"price" binding:"required"`
	InitialQuantity   int       `json:"initial_quantity" binding:"required"`
	AvailableQuantity int       `json:"available_quantity" binding:"required"`
	Location          []float64 `json:"location" binding:"required"`
	Created           time.Time `json:"created" binding:"required"`
}
