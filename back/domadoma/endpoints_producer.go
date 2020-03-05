package domadoma

import "net/http"

type ProducerEndpoints interface{
	GetProducer() func(w http.ResponseWriter, r *http.Request)

	CreateProducer() func(w http.ResponseWriter, r *http.Request)

	ListProducers() func(w http.ResponseWriter, r *http.Request)

	UpdateProducer() func(w http.ResponseWriter, r *http.Request)

	DeleteProducer() func(w http.ResponseWriter, r *http.Request)

}

func NewProducerEndpoints(producerBase ProducerBase) ProducerEndpoints {
	return &ProducerEndpointsFactory{producerBase: producerBase}
}

type ProducerEndpointsFactory struct{
	producerBase ProducerBase
}