package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fullacc/photoback/back/photo_base"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
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
	var configfile *photo_base.ConfigFile
	if err := json.Unmarshal(data, &configfile); err != nil {
		return err
	}
	file.Close()

	postgrePersonStore, err := photo_base.NewPostgrePersonStore(configfile)
	if err != nil {
		panic(err)
	}

	postgreOperationStore, err := photo_base.NewPostgreOperationStore(configfile)
	if err != nil {
		panic(err)
	}

	postgrePhotoStore, err := photo_base.NewPostgrePhotoStore(configfile)
	if err != nil {
		panic(err)
	}

	postgrepersonendpoints := photo_base.NewEndpointsPersonFactory(postgrePersonStore)
	postgreoperationendpoints := photo_base.NewEndpointsOperationFactory(postgreOperationStore)
	postgrephotoendpoints := photo_base.NewEndpointsPhotoFactory(postgrePhotoStore)

	router := mux.NewRouter()

	router.Methods("GET").Path("/person/{id}").HandlerFunc(postgrepersonendpoints.GetPerson("id"))
	router.Methods("POST").Path("/person/").HandlerFunc(postgrepersonendpoints.CreatePerson())
	router.Methods("GET").Path("/person/").HandlerFunc(postgrepersonendpoints.ListPersons())
	router.Methods("PUT").Path("/person/{id}").HandlerFunc(postgrepersonendpoints.UpdatePerson("id"))
	router.Methods("DELETE").Path("/person/{id}").HandlerFunc(postgrepersonendpoints.DeletePerson("id"))

	router.Methods("GET").Path("/person/{personid}/operation/{id}").HandlerFunc(postgreoperationendpoints.GetOperation("id"))
	router.Methods("POST").Path("/person/{personid}/operation/").HandlerFunc(postgreoperationendpoints.CreateOperation())
	router.Methods("GET").Path("/operation/").HandlerFunc(postgreoperationendpoints.ListOperations())
	router.Methods("GET").Path("/person/{id}/operation").HandlerFunc(postgreoperationendpoints.ListPersonOperations("id"))
	router.Methods("PUT").Path("/person/{personid}/operation/{id}").HandlerFunc(postgreoperationendpoints.UpdateOperation("id"))
	router.Methods("DELETE").Path("/person/{personid}/operation/{id}").HandlerFunc(postgreoperationendpoints.DeleteOperation("personid","id"))

	router.Methods("GET").Path("/person/{personid}/operation/{operationid}/photo/{id}").HandlerFunc(postgrephotoendpoints.GetPhoto("id"))
	router.Methods("POST").Path("/person/{personid}/operation/{operationid}/photo/").HandlerFunc(postgrephotoendpoints.CreatePhoto("personid","operationid"))
	router.Methods("GET").Path("/photo/").HandlerFunc(postgrephotoendpoints.ListPhotos())
	router.Methods("GET").Path("/person/{personid}/photo").HandlerFunc(postgrephotoendpoints.ListPersonPhotos("personid"))
	router.Methods("GET").Path("/person/{personid}/operation/{operationid}/photo/").HandlerFunc(postgrephotoendpoints.ListOperationPhotos("operationid"))
	router.Methods("PUT").Path("/person/{personid}/operation/{operationid}/photo/{id}").HandlerFunc(postgrephotoendpoints.UpdatePhoto("id"))
	router.Methods("DELETE").Path("/person/{personid}/operation/{operationid}/photo/{id}").HandlerFunc(postgrephotoendpoints.DeletePhoto("id"))


	fmt.Println("Server started")
	go func(port string, rtr *mux.Router) {
		http.ListenAndServe("0.0.0.0:" + port, rtr)
	}(configfile.Port, router)

	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt,syscall.SIGTERM)
	go func() {
		<-c
		done <- true
	}()

	<- done
	log.Printf("server shutdown")
	os.Exit(1)

	return nil
}