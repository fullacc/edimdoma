package Feedback

import "time"

type FeedbackBase interface{
	CreateFeedback(feedback *Feedback) (*Feedback, error)

	GetFeedback(id int) (*Feedback, error)

	ListFeedbacks() ([]*Feedback, error)

	UpdateFeedback(id int, feedback *Feedback) (*Feedback, error)

	DeleteFeedback(id int)  error
}

type Feedback struct {
	Id int `json:"id"`
	ProducerId int `json:"producer_id" binding:"required"`
	ConsumerId int `json:"consumer_id" binding:"required"`
	Value int `json:"value" binding:"required"`
	Text string `json:"text" binding:"required"`
	Created time.Time `json:"created" binding:"required"`
	FeedbackId int `json:"deal_id" binding:"required"`
}
