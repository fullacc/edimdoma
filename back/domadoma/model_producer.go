package domadoma

type ProducerBase interface{
	CreateProducer()

	GetProducer()

	ListProducers()

	UpdateProducer()

	DeleteProducer()
}

type Producer struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	Rating float64 `json:"rating"`

}