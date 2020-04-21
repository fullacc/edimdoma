package Offer

import (
	"bufio"
	"fmt"
	"github.com/segmentio/encoding/json"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/fullacc/edimdoma/back/domadoma/Connection"
	"github.com/go-pg/pg"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var offerBase OfferBase
var db *pg.DB


func TestPostgreOfferBase_CreateOffer(t *testing.T) {
	var err error
	file, err := os.Open("./../../../../4ernovik/config.json")
	if err != nil {
		t.Error(err)
	}

	buffer := bufio.NewReader(file)
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		t.Error(err)
	}

	var configfile *domadoma.ConfigFile
	if err := json.Unmarshal(data, &configfile); err != nil {
		t.Error(err)
	}
	_ = file.Close()
	db = Connection.ConnectToPostgre(configfile)
	offerBase, err = NewPostgreOfferBase(db)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:57:57.244503+06:00")
	if err != nil {
		t.Error(err)
	}
	offer:=&Offer{
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"ggwp",
		Created: crtd,
		DealId:1,
		Anon:1,
	}
	offer,err = offerBase.CreateOffer(offer)
	if err != nil {
		t.Error(err)
	}
	offertest:=Offer{
		Id:offer.Id,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"ggwp",
		Created: crtd,
		DealId:1,
		Anon:1,}
	if *offer != offertest {
		t.Error("Not equal")
	}
}

func TestPostgreOfferBase_GetOffer(t *testing.T) {
	var err error
	file, err := os.Open("./../../../../4ernovik/config.json")
	if err != nil {
		t.Error(err)
	}

	buffer := bufio.NewReader(file)
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		t.Error(err)
	}

	var configfile *domadoma.ConfigFile
	if err := json.Unmarshal(data, &configfile); err != nil {
		t.Error(err)
	}
	_ = file.Close()
	db = Connection.ConnectToPostgre(configfile)
	offerBase, err = NewPostgreOfferBase(db)
	if err != nil {
		t.Error(err)
	}
	offer := &Offer{Id: 1}
	offer, err = offerBase.GetOffer(offer)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:59:02.440137+06:00")
	if err != nil {
		t.Error(err)
	}

	offertest:=Offer{
		Id:1,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"very tasty BESH",
		Created: crtd,
		DealId:1,
		Anon:1,}
	if *offer != offertest{
		t.Error("Title is not equal")
	}
}

func TestPostgreOfferBase_ListOffers(t *testing.T) {
	var err error
	file, err := os.Open("./../../../../4ernovik/config.json")
	if err != nil {
		t.Error(err)
	}

	buffer := bufio.NewReader(file)
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		t.Error(err)
	}

	var configfile *domadoma.ConfigFile
	if err := json.Unmarshal(data, &configfile); err != nil {
		t.Error(err)
	}
	_ = file.Close()
	db = Connection.ConnectToPostgre(configfile)
	offerBase, err = NewPostgreOfferBase(db)
	if err != nil {
		t.Error(err)
	}
	offers, err := offerBase.ListOffers()
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:59:02.440137+06:00")
	if err != nil {
		t.Error(err)
	}

	offertest:=Offer{
		Id:offers[0].Id,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"very tasty BESH",
		Created: crtd,
		DealId:1,
		Anon:1,}
	if *offers[0] != offertest{
		fmt.Println(offertest)
		fmt.Println(*offers[0])
		t.Error("Title is not equal")
	}
}

func TestPostgreOfferBase_DeleteOffer(t *testing.T) {
	var err error
	file, err := os.Open("./../../../../4ernovik/config.json")
	if err != nil {
		t.Error(err)
	}

	buffer := bufio.NewReader(file)
	data, err := ioutil.ReadAll(buffer)
	if err != nil {
		t.Error(err)
	}

	var configfile *domadoma.ConfigFile
	if err := json.Unmarshal(data, &configfile); err != nil {
		t.Error(err)
	}
	_ = file.Close()
	db = Connection.ConnectToPostgre(configfile)
	offerBase, err = NewPostgreOfferBase(db)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:57:57.244503+06:00")
	if err != nil {
		t.Error(err)
	}
	offer:=&Offer{
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"very tasty BESH",
		Created: crtd,
		DealId:1,
		Anon:1,}
	offer,err = offerBase.CreateOffer(offer)
	if err != nil {
		t.Error(err)
	}
	offerBase.DeleteOffer(offer.Id)
	_, err = offerBase.GetOffer(offer)
	if err == nil {
		t.Errorf("Not deleted")
	}
}
