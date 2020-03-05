package domadoma

import "time"

type FeedbackBase interface{
	CreateFeedback()

	GetFeedback()

	ListFeedback()

	UpdateFeedback()

	DeleteFeedback()
}

type Feedback struct {
	Id int `json:"id"`
	ProducerId int `json:"producer_id"`
	ConsumerId int `json:"consumer_id"`
	Value int `json:"value"`
	Text string `json:"text"`
	Created time.Time `json:"created"`
	FeedbackId int `json:"deal_id"`
}
