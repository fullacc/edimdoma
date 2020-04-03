package Feedback

import "time"

type FeedbackBase interface {
	CreateFeedback(feedback *Feedback) (*Feedback, error)

	GetFeedback(feedback *Feedback) (*Feedback, error)

	ListFeedbacks() ([]*Feedback, error)

	ListProducerFeedbacks(id int) ([]*Feedback, error)

	UpdateFeedback(feedback *Feedback) (*Feedback, error)

	DeleteFeedback(id int) error
}

type Feedback struct {
	Id         int       `json:"id"`
	ProducerId int       `json:"producer_id"`
	ConsumerId int       `json:"consumer_id"`
	Value      int       `json:"value" binding:"required"`
	Text       string    `json:"text" binding:"required"`
	Created    time.Time `json:"created"`
	DealId     int       `json:"deal_id"`
}
