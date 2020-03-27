package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"./domadoma"
	"./domadoma/Deal"
	"./domadoma/Feedback"
	"./domadoma/Offer"
	"./domadoma/OfferLog"
	"./domadoma/Request"
	"./domadoma/User"
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


	postgreDealBase, err := Deal.NewPostgreDealBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreFeedbackBase, err := Feedback.NewPostgreFeedbackBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreOfferBase, err := Offer.NewPostgreOfferBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreOfferLogBase, err := OfferLog.NewPostgreOfferLogBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreRequestBase, err := Request.NewPostgreRequestBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreUserBase, err := User.NewPostgreUserBase(configfile)
	if err != nil {
		panic(err)
	}

	postgreDealEndpoints := Deal.NewDealEndpoints(postgreDealBase)

	postgreFeedbackEndpoints := Feedback.NewFeedbackEndpoints(postgreFeedbackBase)

	postgreOfferEndpoints := Offer.NewOfferEndpoints(postgreOfferBase)

	postgreOfferLogEndpoints := OfferLog.NewOfferLogEndpoints(postgreOfferLogBase)

	postgreRequestEndpoints := Request.NewRequestEndpoints(postgreRequestBase)

	postgreUserEndpoints := User.NewUserEndpoints(postgreUserBase)

	router := gin.Default()

	api := router.Group("api")
	{

		api.GET("deal",postgreDealEndpoints.ListDeals())
		api.GET("feedback",postgreFeedbackEndpoints.ListFeedbacks())
		api.GET("offer",postgreOfferEndpoints.ListOffers())
		api.GET("offer_log",postgreOfferLogEndpoints.ListOfferLogs())
		api.GET("request",postgreRequestEndpoints.ListRequests())
		api.GET("user",postgreUserEndpoints.ListUsers())



		deals := api.Group("deal")
		{
			deals.GET(":dealid",postgreDealEndpoints.GetDeal())
			deals.DELETE(":dealid",postgreDealEndpoints.DeleteDeal())
			deals.PUT(":dealid",postgreDealEndpoints.UpdateDeal())
			deals.PATCH(":dealid/complete",postgreDealEndpoints.CompleteDeal())
		}

		consumers := api.Group("consumer")
		{
			requests := consumers.Group(":consumerid/request")
			{
				requests.POST("",postgreRequestEndpoints.CreateRequest())
				requests.GET("",postgreRequestEndpoints.ListRequests())
				requests.GET(":requestid",postgreRequestEndpoints.GetRequest())
				requests.PUT(":requestid",postgreRequestEndpoints.UpdateRequest())
				requests.DELETE(":requestid",postgreRequestEndpoints.DeleteRequest())
			}
			consumers.POST(":consumerid/offer/:offerid",postgreDealEndpoints.CreateDeal())//needs quantity and price field
			consumers.GET(":consumerid/deal",postgreDealEndpoints.ListDeals())//requires onlyactive query

			consumers.POST(":consumerid/deal/:dealid/feedback",postgreFeedbackEndpoints.CreateFeedback())//only value and text
			consumers.GET(":consumerid/feedback/:feedbackid",postgreFeedbackEndpoints.GetFeedback())
			consumers.PUT(":consumerid/feedback/:feedbackid",postgreFeedbackEndpoints.UpdateFeedback())
			consumers.DELETE(":consumerid/feedback/:feedbackid",postgreFeedbackEndpoints.DeleteFeedback())
		}

		producers := api.Group("producer")
		{
			offers := producers.Group(":producerid/offer")
			{
				offers.POST("",postgreOfferEndpoints.CreateOffer())
				offers.GET("",postgreOfferEndpoints.ListOffers())
				offers.GET(":offerid",postgreOfferEndpoints.GetOffer())
				offers.PUT(":offerid",postgreOfferEndpoints.UpdateOffer())
				offers.DELETE(":offerid",postgreOfferEndpoints.DeleteOffer())
			}
			producers.POST(":producerid/request/:requestid",postgreDealEndpoints.CreateDeal())//needs price only
			producers.GET(":producerid/deal",postgreDealEndpoints.ListDeals())//requires onlyactive query

			producers.GET(":producerid/feedback",postgreFeedbackEndpoints.ListFeedbacks())
			offerlogs := producers.Group(":producerid/offerlog")
			{
				offerlogs.GET("",postgreOfferLogEndpoints.ListOfferLogs())
				offerlogs.GET(":offerlogid",postgreOfferLogEndpoints.GetOfferLog())
				offerlogs.DELETE(":offerlogid",postgreOfferLogEndpoints.DeleteOfferLog())
			}
		}

		users := api.Group("user")
		{
			users.POST("registration",postgreUserEndpoints.CreateUser())
			users.POST("login",postgreUserEndpoints.LoginUser())
			users.PATCH("logout",postgreUserEndpoints.LogoutUser())
			users.GET(":userid",postgreUserEndpoints.GetUser())
			users.PUT(":userid",postgreUserEndpoints.UpdateUser())
			users.DELETE(":userid",postgreUserEndpoints.DeleteUser())
			users.PUT(":userid/changepassword",postgreUserEndpoints.ChangePassword())
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