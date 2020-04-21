package main

import (
	"bufio"
	"fmt"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/Connection"
	"github.com/fullacc/edimdoma/back/domadoma/Deal"
	"github.com/fullacc/edimdoma/back/domadoma/Endpoints"
	"github.com/fullacc/edimdoma/back/domadoma/Feedback"
	"github.com/fullacc/edimdoma/back/domadoma/Offer"
	"github.com/fullacc/edimdoma/back/domadoma/OfferLog"
	"github.com/fullacc/edimdoma/back/domadoma/Rabbit"
	"github.com/fullacc/edimdoma/back/domadoma/Request"
	"github.com/fullacc/edimdoma/back/domadoma/SMS"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
	"github.com/segmentio/encoding/json"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	config = ""
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "config /filepath",
				Required:    true,
				Destination: &config,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "launch",
				Aliases: []string{"l"},
				Usage:   "launch",
				Action:  run,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(g *cli.Context) error {
	file, err := os.Open(config)
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
	_ = file.Close()

	router ,err := setupRouter(configfile)
	if err != nil {
		return err
	}

	go func(port string, rtr *gin.Engine) {
		rtr.Run("0.0.0.0:" + port)
	}(configfile.ApiPort, router)
	fmt.Println("Server started")


	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		done <- true
	}()

	<-done
	log.Printf("server shutdown")
	os.Exit(1)
	return nil
}

func setupRouter(configfile *domadoma.ConfigFile)  (*gin.Engine,error) {

	var db *pg.DB
	db = Connection.ConnectToPostgre(configfile)

	var red *redis.Client
	red = Connection.ConnectToRedis(configfile)

	postgreDealBase, err := Deal.NewPostgreDealBase(db)
	if err != nil {
		return nil, err
	}

	postgreFeedbackBase, err := Feedback.NewPostgreFeedbackBase(db)
	if err != nil {
		return nil, err
	}

	postgreOfferBase, err := Offer.NewPostgreOfferBase(db)
	if err != nil {
		return nil, err
	}

	postgreOfferLogBase, err := OfferLog.NewPostgreOfferLogBase(db)
	if err != nil {
		return nil, err
	}

	postgreRequestBase, err := Request.NewPostgreRequestBase(db)
	if err != nil {
		return nil, err
	}

	postgreUserBase, err := User.NewPostgreUserBase(db)
	if err != nil {
		return nil, err
	}

	redisAuthorizationBase, err := Authorization.NewRedisAuthorizationBase(red)
	if err != nil {
		return nil, err
	}

	smsBase, err := SMS.NewSMSBase(configfile)
	if err != nil {
		return nil, err
	}

	rabbitBase, err := Rabbit.NewRabbitMQRabbitBase(configfile,postgreOfferBase,postgreOfferLogBase)

	postgreDealEndpoints := Endpoints.NewDealEndpoints(postgreDealBase, redisAuthorizationBase, postgreOfferBase, postgreOfferLogBase, postgreRequestBase, postgreUserBase)

	postgreFeedbackEndpoints := Endpoints.NewFeedbackEndpoints(postgreFeedbackBase, redisAuthorizationBase, postgreDealBase, postgreUserBase)

	postgreOfferEndpoints := Endpoints.NewOfferEndpoints(postgreOfferBase, redisAuthorizationBase, postgreUserBase, rabbitBase)

	postgreOfferLogEndpoints := Endpoints.NewOfferLogEndpoints(postgreOfferLogBase, redisAuthorizationBase)

	postgreRequestEndpoints := Endpoints.NewRequestEndpoints(postgreRequestBase, redisAuthorizationBase, postgreUserBase)

	postgreUserEndpoints := Endpoints.NewUserEndpoints(postgreUserBase, redisAuthorizationBase)

	redisAuthorizationEndpoints := Endpoints.NewAuthorizationEndpoints(redisAuthorizationBase, smsBase, postgreUserBase)

	go func() {
		rabbitBase.ConsumeRabbit()
	}()


	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORS())


	api := router.Group("api")
	{

		api.GET("deal", postgreDealEndpoints.ListDeals())
		api.GET("feedback", postgreFeedbackEndpoints.ListFeedbacks())
		api.GET("offer", postgreOfferEndpoints.ListOffers())
		api.GET("offer_log", postgreOfferLogEndpoints.ListOfferLogs())
		api.GET("request", postgreRequestEndpoints.ListRequests())
		api.GET("request/:requestid", postgreRequestEndpoints.GetRequest())
		api.GET("user", postgreUserEndpoints.ListUsers())
		api.POST("user", postgreUserEndpoints.CreateUser())

		deals := api.Group("deal")
		{
			deals.GET(":dealid", postgreDealEndpoints.GetDeal())
			deals.DELETE(":dealid", postgreDealEndpoints.DeleteDeal())
			deals.PUT(":dealid", postgreDealEndpoints.UpdateDeal())
			deals.PATCH(":dealid/complete", postgreDealEndpoints.CompleteDeal())
		}

		consumers := api.Group("consumer")
		{
			requests := consumers.Group(":consumerid/request")
			{
				requests.POST("", postgreRequestEndpoints.CreateRequest())
				requests.GET("", postgreRequestEndpoints.ListRequests())

				requests.PUT(":requestid", postgreRequestEndpoints.UpdateRequest())
				requests.DELETE(":requestid", postgreRequestEndpoints.DeleteRequest())
			}
			consumers.POST(":consumerid/offer/:offerid", postgreDealEndpoints.CreateDeal()) //needs quantity and price field
			consumers.GET(":consumerid/deal", postgreDealEndpoints.ListDeals())             //requires onlyactive query

			consumers.POST(":consumerid/deal/:dealid/feedback", postgreFeedbackEndpoints.CreateFeedback()) //only value and text
			consumers.GET(":consumerid/feedback/:feedbackid", postgreFeedbackEndpoints.GetFeedback())
			consumers.PUT(":consumerid/feedback/:feedbackid", postgreFeedbackEndpoints.UpdateFeedback())
			consumers.DELETE(":consumerid/feedback/:feedbackid", postgreFeedbackEndpoints.DeleteFeedback())
		}

		producers := api.Group("producer")
		{
			offers := producers.Group(":producerid/offer")
			{
				offers.POST("", postgreOfferEndpoints.CreateOffer())
				offers.GET("", postgreOfferEndpoints.ListOffers())
				offers.GET(":offerid", postgreOfferEndpoints.GetOffer())
				offers.PUT(":offerid", postgreOfferEndpoints.UpdateOffer())
				offers.DELETE(":offerid", postgreOfferEndpoints.DeleteOffer())
			}
			producers.POST(":producerid/request/:requestid", postgreDealEndpoints.CreateDeal()) //needs price only
			producers.GET(":producerid/deal", postgreDealEndpoints.ListDeals())                 //requires onlyactive query

			producers.GET(":producerid/feedback", postgreFeedbackEndpoints.ListFeedbacks())
			offerlogs := producers.Group(":producerid/offerlog")
			{
				offerlogs.GET("", postgreOfferLogEndpoints.ListOfferLogs())
				offerlogs.GET(":offerlogid", postgreOfferLogEndpoints.GetOfferLog())
				offerlogs.DELETE(":offerlogid", postgreOfferLogEndpoints.DeleteOfferLog())
			}
		}

		users := api.Group("user")
		{
			users.GET(":userid", postgreUserEndpoints.GetUser())
			users.PUT(":userid", postgreUserEndpoints.UpdateUser())
			users.DELETE(":userid", postgreUserEndpoints.DeleteUser())
			users.POST("registration", redisAuthorizationEndpoints.RegisterUser())
			users.POST("checkphone", redisAuthorizationEndpoints.CheckPhone())
			users.POST("checkcode", redisAuthorizationEndpoints.CheckCode())
			users.POST("login", redisAuthorizationEndpoints.LoginUser())
			users.PATCH("logout", redisAuthorizationEndpoints.LogoutUser())
			users.PUT(":userid/changepassword", redisAuthorizationEndpoints.ChangePassword())
		}
	}

	return router, nil
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}