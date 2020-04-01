package OfferLog

import "time"

type OfferLogBase interface{
	CreateOfferLog(offerLog *OfferLog) (*OfferLog, error)

	GetOfferLog(offerLog *OfferLog) (*OfferLog, error)

	ListOfferLogs() ([]*OfferLog, error)

	ListProducerOfferLogs(id int) ([]*OfferLog, error)

	UpdateOfferLog(offerLog *OfferLog) (*OfferLog, error)

	DeleteOfferLog(id int)  error
}

type OfferLog struct {
	Id int `json:"id"`
	ProducerId int `json:"producer_id" binding:"required"`
	Food string `json:"name" binding:"required"`
	Price int `json:"price" binding:"required"`
	InitialQuantity int `json:"initial_quantity" binding:"required"`
	AvailableQuantity int `json:"available_quantity" binding:"required"`
	Location []float64 `json:"location" binding:"required"`
	Created time.Time `json:"created" binding:"required"`
}

