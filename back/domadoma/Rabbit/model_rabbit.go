package Rabbit

type RabbitBase interface {
	CreateRabbit(id int) (int, error)

	ConsumeRabbit() error
}

const DelayMilliseconds string = "28800000"