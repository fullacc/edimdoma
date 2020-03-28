package Request

import "time"

type RequestBase interface{
	CreateRequest(request *Request) (*Request, error)

	GetRequest(id int) (*Request, error)

	ListRequests() ([]*Request, error)

	ListConsumerRequests(id int) ([]*Request, error)

	UpdateRequest(request *Request) (*Request, error)

	DeleteRequest(id int)  error
}

type Request struct {
	Id int `json:"id"`
	ConsumerId int `json:"consumer_id"`
	Food string `json:"name" binding:"required"`
	Price int `json:"price" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
	Location []float64 `json:"location" binding:"required"`
	Created time.Time `json:"created"`
}

