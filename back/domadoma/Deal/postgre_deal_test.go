package Deal

import (
	"bufio"
	"github.com/segmentio/encoding/json"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/fullacc/edimdoma/back/domadoma/Connection"
	"github.com/go-pg/pg"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var dealBase DealBase
var db *pg.DB


func TestPostgreDealBase_CreateDeal(t *testing.T) {
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
	dealBase, err = NewPostgreDealBase(db)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:57:57.244503+06:00")
	if err != nil {
		t.Error(err)
	}
	zr, err := time.Parse("2006-01-02T15:04:05Z07:00","0001-01-01T00:00:00Z")
	if err != nil {
		t.Error(err)
	}
	deal:=&Deal{
		FoodName:"test",
		Type:1,
		Myaso:1,
		Halal:1,
		Vegan:1,
		Spicy:3,
		Description:"testtext",
		Price:69,
		Quantity:3,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		ProducerName:"stella",
		Created: crtd,
		Finished: zr,
		Complete:"false"}
	deal,err = dealBase.CreateDeal(deal)
	if err != nil {
		t.Error(err)
	}
	dealtest:=Deal{
		Id:deal.Id,
		FoodName:"test",
		Type:1,
		Myaso:1,
		Halal:1,
		Vegan:1,
		Spicy:3,
		Description:"testtext",
		Price:69,
		Quantity:3,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		ProducerName:"stella",
		Created: crtd,
		Finished: zr,
		Complete:"false"}
	if *deal != dealtest {
		t.Error("Not equal")
	}
}

func TestPostgreDealBase_GetDeal(t *testing.T) {
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
	dealBase, err = NewPostgreDealBase(db)
	if err != nil {
		t.Error(err)
	}
	deal := &Deal{Id: 1}
	deal, err = dealBase.GetDeal(deal)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:57:57.244503+06:00")
	if err != nil {
		t.Error(err)
	}
	zr, err := time.Parse("2006-01-02T15:04:05Z07:00","0001-01-01T00:00:00Z")
	if err != nil {
		t.Error(err)
	}
	dealtest:=Deal{
		Id:1,
		FoodName:"lasagna",
		Type:0,
		Myaso:0,
		Halal:0,
		Vegan:0,
		Spicy:2,
		Description:"",
		Price:12312312,
		Quantity:1,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		ProducerName:"stella",
		Created: crtd,
		Finished: zr,
		Complete:"false"}
	if *deal != dealtest{
		t.Error("Title is not equal")
	}
}

func TestPostgreDealBase_ListDeals(t *testing.T) {
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
	dealBase, err = NewPostgreDealBase(db)
	if err != nil {
		t.Error(err)
	}
	deals, err := dealBase.ListDeals()
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:57:57.244503+06:00")
	if err != nil {
		t.Error(err)
	}
	zr, err := time.Parse("2006-01-02T15:04:05Z07:00","0001-01-01T00:00:00Z")
	if err != nil {
		t.Error(err)
	}
	dealtest:=Deal{
		Id:deals[0].Id,
		FoodName:"lasagna",
		Type:0,
		Myaso:0,
		Halal:0,
		Vegan:0,
		Spicy:2,
		Description:"",
		Price:12312312,
		Quantity:1,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		ProducerName:"stella",
		Created: crtd,
		Finished: zr,
		Complete:"false"}
	if *deals[0] != dealtest{
		//fmt.Println(dealtest)
		//fmt.Println(*deals[1])
		t.Error("Title is not equal")
	}
}

func TestPostgreDealBase_DeleteDeal(t *testing.T) {
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
	dealBase, err = NewPostgreDealBase(db)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:57:57.244503+06:00")
	if err != nil {
		t.Error(err)
	}
	zr, err := time.Parse("2006-01-02T15:04:05Z07:00","0001-01-01T00:00:00Z")
	if err != nil {
		t.Error(err)
	}
	deal:=&Deal{
		FoodName:"test",
		Type:1,
		Myaso:1,
		Halal:1,
		Vegan:1,
		Spicy:3,
		Description:"testtext",
		Price:69,
		Quantity:3,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		ProducerName:"stella",
		Created: crtd,
		Finished: zr,
		Complete:"false"}
	deal,err = dealBase.CreateDeal(deal)
	if err != nil {
		t.Error(err)
	}
	dealBase.DeleteDeal(deal.Id)
	_, err = dealBase.GetDeal(deal)
	if err == nil {
		t.Errorf("Not deleted")
	}
}