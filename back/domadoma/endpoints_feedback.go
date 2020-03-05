package domadoma

import "net/http"

type FeedbackEndpoints interface{
	GetFeedback() func(w http.ResponseWriter, r *http.Request)

	CreateFeedback() func(w http.ResponseWriter, r *http.Request)

	ListFeedbacks() func(w http.ResponseWriter, r *http.Request)

	UpdateFeedback() func(w http.ResponseWriter, r *http.Request)

	DeleteFeedback() func(w http.ResponseWriter, r *http.Request)

}

func NewFeedbackEndpoints(feedbackBase FeedbackBase) FeedbackEndpoints {
	return &FeedbackEndpointsFactory{feedbackBase: feedbackBase}
}

type FeedbackEndpointsFactory struct{
	feedbackBase FeedbackBase
}