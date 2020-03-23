package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	port=""
	config="./config.json"
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "config /filepath",
			Destination: &config,
		},
	}
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "launch",
				Aliases: []string{"l"},
				Usage:   "launch",
				Action:  run,
			},
		},
	}
	app.Flags=flags
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err:=LaunchServer(config);err!=nil {
		return err
	}
	return nil
}

func LaunchServer(configpath string) error{
	file, err := os.Open(configpath)
	if err != nil {
		return err
	}
	buffer := bufio.NewReader(file)
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		return err
	}
	var configfile *domadoma.ConfigFile
	if err := json.Unmarshal(data, &configfile); err != nil {
		return err
	}
	file.Close()


	postgreDealBase, err := domadoma.NewPostgreDealBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreFeedbackBase, err := domadoma.NewPostgreFeedbackBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreOfferBase, err := domadoma.NewPostgreOfferBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreOfferLogBase, err := domadoma.NewPostgreOfferLogBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreRequestBase, err := domadoma.NewPostgreRequestBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreUserBase, err := domadoma.NewPostgreUserBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreDealEndpoints := domadoma.NewDealEndpoints(postgreDealBase)

	postgreFeedbackEndpoints := domadoma.NewFeedbackEndpoints(postgreFeedbackBase)

	postgreOfferEndpoints := domadoma.NewOfferEndpoints(postgreOfferBase)

	postgreOfferLogEndpoints := domadoma.NewOfferLogEndpoints(postgreOfferLogBase)

	postgreRequestEndpoints := domadoma.NewRequestEndpoints(postgreRequestBase)

	postgreUserEndpoints := domadoma.NewUserEndpoints(postgreUserBase)

	router := gin.Default()

	api := router.Group("api")
	{

		api.GET("deals",postgreDealEndpoints.ListDeals())
		api.GET("feedbacks",postgreFeedbackEndpoints.ListFeedbacks())
		api.GET("offers",postgreOfferEndpoints.ListOffers())
		api.GET("offer_logs",postgreOfferLogEndpoints.ListOfferLogs())
		api.GET("requests",postgreRequestEndpoints.ListRequests())
		api.GET("users",postgreUserEndpoints.ListUsers())

		deals := api.Group("deal")
		{
			deals.GET("",postgreDealEndpoints.ListDeals())
			deals.GET(":dealid",postgreDealEndpoints.GetDeal())
			deals.DELETE(":dealid",postgreDealEndpoints.DeleteDeal())
			deals.PUT(":dealid",postgreDealEndpoints.UpdateDeal())
		}

		consumers := api.Group("consumer")
		{
			requests := consumers.Group(":consumerid/request")
			{
				requests.POST("",postgreRequestEndpoints.CreateRequest())
				requests.GET("",postgreRequestEndpoints.ListConsumerRequests())
				requests.GET(":requestid",postgreRequestEndpoints.GetRequest())
				requests.PUT(":requestid",postgreRequestEndpoints.UpdateRequest())
				requests.DELETE(":requestid",postgreRequestEndpoints.DeleteRequest())
			}
			consumers.POST(":consumerid/offer/:offerid",postgreDealEndpoints.CreateDeal())
			consumers.GET(":consumerid/deal",postgreDealEndpoints.ListDeals())
			consumers.GET(":consumerid/deal/:dealid",postgreDealEndpoints.GetDeal())
			consumers.POST(":consumerid/deal/:dealid/feedback",postgreFeedbackEndpoints.CreateFeedback())
		}

		producers := api.Group("producer")
		{
			offers := producers.Group(":producerid/offer")
			{
				offers.POST("",postgreOfferEndpoints.CreateOffer())
				offers.GET("",postgreOfferEndpoints.ListProducerOffers())
				offers.GET(":offerid",postgreOfferEndpoints.GetOffer())
				offers.PUT(":offerid",postgreOfferEndpoints.UpdateOffer())
				offers.DELETE(":offerid",postgreOfferEndpoints.DeleteOffer())
			}
			producers.POST(":producerid/request/:requestid",postgreDealEndpoints.CreateDeal())
			producers.GET(":producerid/deal",postgreDealEndpoints.ListDeals())
			producers.GET(":producerid/deal/:dealid",postgreDealEndpoints.GetDeal())
			offerlogs := producers.Group(":producerid/offerlog")
			{
				offerlogs.GET("",postgreOfferLogEndpoints.ListProducerOfferLogs())
				offerlogs.GET(":offerlogid",postgreOfferLogEndpoints.GetOfferLog())
				offerlogs.DELETE(":offerlogid",postgreOfferLogEndpoints.DeleteOfferLog())
			}
		}

		users := api.Group("user")
		{
			users.POST("registration",postgreUserEndpoints.CreateUser())
			users.POST("login",postgreUserEndpoints.LoginUser())
			users.GET(":userid",postgreUserEndpoints.GetUser())
			users.PUT(":userid",postgreUserEndpoints.UpdateUser())
			users.DELETE(":userid",postgreUserEndpoints.DeleteUser())
			feedbacks := api.Group(":userid/feedback")
			{
				feedbacks.GET("",postgreFeedbackEndpoints.ListFeedbacks())
				feedbacks.GET(":feedbackid",postgreFeedbackEndpoints.GetFeedback())
				feedbacks.PUT(":feedbackid",postgreFeedbackEndpoints.UpdateFeedback())
				feedbacks.DELETE(":feedbackid",postgreFeedbackEndpoints.DeleteFeedback())
			}
		}
	}

	fmt.Println("Server started")
	go func(port string, rtr *gin.Engine) {
		rtr.Run("0.0.0.0:" + port)
	}(configfile.Port, router)

	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt,syscall.SIGTERM)
	go func(){
		<-c
		done <- true
	}()

	<- done
	log.Printf("server shutdown")
	os.Exit(1)

	return nil
}