package Deal

import "time"

type DealBase interface {
	CreateDeal(deal *Deal) (*Deal, error)

	GetDeal(deal *Deal) (*Deal, error)

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
	Id           int       `json:"id"`
	FoodName     string    `json:"food_name"`
	Type         int       `json:"type"`
	Myaso        int       `json:"myaso"`
	Halal        int       `json:"halal"`
	Vegan        int       `json:"vegan"`
	Spicy        int       `json:"spicy"`
	Description  string    `json:"description"`
	Price        int       `json:"price" binding:"required"`
	Quantity     int       `json:"quantity"`
	Comment      int       `json:"comment"`
	ConsumerId   int       `json:"consumer_id"`
	ProducerId   int       `json:"producer_id"`
	ConsumerName string    `json:"consumer_name"`
	ProducerName string    `json:"producer_name"`
	Created      time.Time `json:"created"`
	Finished     time.Time `json:"finished,omitempty"`
	Complete     string    `json:"complete"`
}
