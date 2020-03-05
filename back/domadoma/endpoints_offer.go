package domadoma

import "net/http"

type OfferEndpoints interface{
	GetOffer() func(w http.ResponseWriter, r *http.Request)

	CreateOffer() func(w http.ResponseWriter, r *http.Request)

	ListOffers() func(w http.ResponseWriter, r *http.Request)

	UpdateOffer() func(w http.ResponseWriter, r *http.Request)

	DeleteOffer() func(w http.ResponseWriter, r *http.Request)

}

func NewOfferEndpoints(offerBase OfferBase) OfferEndpoints {
	return &OfferEndpointsFactory{offerBase: offerBase}
}

type OfferEndpointsFactory struct{
	offerBase OfferBase
}