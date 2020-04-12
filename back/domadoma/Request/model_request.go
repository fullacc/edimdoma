package Request

import "time"

type RequestBase interface {
	CreateRequest(request *Request) (*Request, error)

	GetRequest(request *Request) (*Request, error)

	ListRequests() ([]*Request, error)

	ListConsumerRequests(id int) ([]*Request, error)

	UpdateRequest(request *Request) (*Request, error)

	DeleteRequest(id int) error
}

type Request struct {
	Id           int       `json:"id"`
	ConsumerId   int       `json:"consumer_id"`
	ConsumerName string    `json:"consumer_name"`
	FoodName     string    `json:"food_name" binding:"required"`
	Type         int       `json:"type"`
	Myaso        int       `json:"myaso"`
	Halal        int       `json:"halal"`
	Vegan        int       `json:"vegan"`
	Spicy        int       `json:"spicy"`
	Description  string    `json:"description"`
	Price        int       `json:"price" binding:"required"`
	Quantity     int       `json:"quantity" binding:"required"`
	Location     []float64 `json:"location" binding:"required"`
	Created      time.Time `json:"created"`
}
