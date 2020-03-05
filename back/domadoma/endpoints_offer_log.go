package domadoma

import "net/http"

type OfferLogEndpoints interface{
	GetOfferLog() func(w http.ResponseWriter, r *http.Request)

	CreateOfferLog() func(w http.ResponseWriter, r *http.Request)

	ListOfferLogs() func(w http.ResponseWriter, r *http.Request)

	UpdateOfferLog() func(w http.ResponseWriter, r *http.Request)

	DeleteOfferLog() func(w http.ResponseWriter, r *http.Request)

}

func NewOfferLogEndpoints(offerLogBase OfferLogBase) OfferLogEndpoints {
	return &OfferLogEndpointsFactory{offerLogBase: offerLogBase}
}

type OfferLogEndpointsFactory struct{
	offerLogBase OfferLogBase
}