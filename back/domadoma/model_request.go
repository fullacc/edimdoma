package domadoma

import "time"

type RequestBase interface{
	CreateRequest()

	GetRequest()

	ListRequests()

	UpdateRequest()

	DeleteRequest()
}

type Request struct {
	Id int `json:"id"`
	ConsumerId int `json:"consumer_id"`
	Name string `json:"name"`
	Price int `json:"price"`
	Quantity int `json:"quantity"`
	Location []float64 `json:"location"`
	Created time.Time `json:"created"`
}

