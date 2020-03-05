package domadoma

import "net/http"

type ConsumerEndpoints interface{
	GetConsumer() func(w http.ResponseWriter, r *http.Request)

	CreateConsumer() func(w http.ResponseWriter, r *http.Request)

	ListConsumers() func(w http.ResponseWriter, r *http.Request)

	UpdateConsumer() func(w http.ResponseWriter, r *http.Request)

	DeleteConsumer() func(w http.ResponseWriter, r *http.Request)

}

func NewConsumerEndpoints(consumerBase ConsumerBase) ConsumerEndpoints {
	return &ConsumerEndpointsFactory{consumerBase: consumerBase}
}

type ConsumerEndpointsFactory struct{
	consumerBase ConsumerBase
}