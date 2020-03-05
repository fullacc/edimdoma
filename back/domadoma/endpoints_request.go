package domadoma

import "net/http"

type RequestEndpoints interface{
	GetRequest() func(w http.ResponseWriter, r *http.Request)

	CreateRequest() func(w http.ResponseWriter, r *http.Request)

	ListRequests() func(w http.ResponseWriter, r *http.Request)

	UpdateRequest() func(w http.ResponseWriter, r *http.Request)

	DeleteRequest() func(w http.ResponseWriter, r *http.Request)

}

func NewRequestEndpoints(requestBase RequestBase) RequestEndpoints {
	return &RequestEndpointsFactory{requestBase: requestBase}
}

type RequestEndpointsFactory struct{
	requestBase RequestBase
}