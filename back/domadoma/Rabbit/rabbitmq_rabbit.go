package Rabbit

import (
	"errors"
	"fmt"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/fullacc/edimdoma/back/domadoma/Offer"
	"github.com/fullacc/edimdoma/back/domadoma/OfferLog"
	"github.com/go-pg/pg"
	"github.com/streadway/amqp"

	"strconv"
)

func NewRabbitMQRabbitBase(configfile *domadoma.ConfigFile,offerBase Offer.OfferBase,offerLogBase OfferLog.OfferLogBase) (RabbitBase, error) {

	rmq, err := amqp.Dial("amqp://"+configfile.RMQLogin+":"+configfile.RMQPassword+"@localhost:"+configfile.RMQPort+"/")
	if err != nil {
		return nil,err
	}

	return &rabbitMQRabbitBase{rmq: rmq,offerLogBase: offerLogBase,offerBase: offerBase}, nil
}

type rabbitMQRabbitBase struct {
	rmq *amqp.Connection
	offerBase         Offer.OfferBase
	offerLogBase OfferLog.OfferLogBase
}

func (r rabbitMQRabbitBase) CreateRabbit(id int) (int, error) {

	ch, err := r.rmq.Channel()
	if err != nil{
		return 0, err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare("exchange","direct",true,false,false,false,nil)
	if err != nil{
		return 0, err
	}
	err = ch.ExchangeDeclare("exchange_dlx","fanout",true,false,false,false,nil)
	if err != nil{
		return 0, err
	}

	_, err = ch.QueueDeclare(
		"queue_dlx", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{"x-dead-letter-exchange":"exchange"},
	)
	if err != nil{
		return 0, err
	}
	err = ch.QueueBind("queue_dlx","key","exchange_dlx",false,nil)
	if err != nil{
		return 0, err
	}

	body := strconv.Itoa(id)
	err = ch.Publish(
		"exchange_dlx",     // exchange
		"key", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
			Expiration: DelayMilliseconds,
		})
	if err != nil{
		return 0, err
	}
	return id,nil
}

func (r rabbitMQRabbitBase) ConsumeRabbit() error{
	ch, err := r.rmq.Channel()
	if err != nil{
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare("exchange","direct",true,false,false,false,nil)
	if err != nil{
		return err
	}

	q, err := ch.QueueDeclare(
		"queue1", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil{
		return err
	}

	err = ch.QueueBind(q.Name,"key","exchange",false,nil)
	if err != nil{
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil{
		return err
	}


	forever := make(chan bool)

	go func() {
		for d := range msgs {
			id, _ := strconv.Atoi(string(d.Body))
			offer := &Offer.Offer{Id:id}
			offer, _ = r.offerBase.GetOffer(offer)
			if err != nil && errors.Is(err, pg.ErrNoRows){
				continue
			}
			offerlog := OfferLog.OfferLog(*offer)
			_, _ = r.offerLogBase.CreateOfferLog(&offerlog)
			_ = r.offerBase.DeleteOffer(offer.Id)
			fmt.Println("Removed ", offer.Id)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
