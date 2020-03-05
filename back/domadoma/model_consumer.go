package domadoma

type ConsumerBase interface{
	Create()

	Get()

	List()

	Update()

	Delete()
}

type Consumer struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
}
