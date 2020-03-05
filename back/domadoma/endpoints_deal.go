package domadoma

import "net/http"

type DealEndpoints interface{
	GetDeal() func(w http.ResponseWriter, r *http.Request)

	CreateDeal() func(w http.ResponseWriter, r *http.Request)

	ListDeals() func(w http.ResponseWriter, r *http.Request)

	UpdateDeal() func(w http.ResponseWriter, r *http.Request)

	DeleteDeal() func(w http.ResponseWriter, r *http.Request)

}

func NewDealEndpoints(dealBase DealBase) DealEndpoints {
	return &DealEndpointsFactory{dealBase: dealBase}
}

type DealEndpointsFactory struct{
	dealBase DealBase
}