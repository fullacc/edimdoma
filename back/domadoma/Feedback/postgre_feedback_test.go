package Feedback

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

var feedbackBase FeedbackBase
var db *pg.DB


func TestPostgreFeedbackBase_CreateFeedback(t *testing.T) {
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
	feedbackBase, err = NewPostgreFeedbackBase(db)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:57:57.244503+06:00")
	if err != nil {
		t.Error(err)
	}
	feedback:=&Feedback{
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"ggwp",
		Created: crtd,
		DealId:1,
		Anon:1,
	}
	feedback,err = feedbackBase.CreateFeedback(feedback)
	if err != nil {
		t.Error(err)
	}
	feedbacktest:=Feedback{
		Id:feedback.Id,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"ggwp",
		Created: crtd,
		DealId:1,
		Anon:1,}
	if *feedback != feedbacktest {
		t.Error("Not equal")
	}
}

func TestPostgreFeedbackBase_GetFeedback(t *testing.T) {
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
	feedbackBase, err = NewPostgreFeedbackBase(db)
	if err != nil {
		t.Error(err)
	}
	feedback := &Feedback{Id: 1}
	feedback, err = feedbackBase.GetFeedback(feedback)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:59:02.440137+06:00")
	if err != nil {
		t.Error(err)
	}

	feedbacktest:=Feedback{
		Id:1,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"very tasty BESH",
		Created: crtd,
		DealId:1,
		Anon:1,}
	if *feedback != feedbacktest{
		t.Error("Title is not equal")
	}
}

func TestPostgreFeedbackBase_ListFeedbacks(t *testing.T) {
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
	feedbackBase, err = NewPostgreFeedbackBase(db)
	if err != nil {
		t.Error(err)
	}
	feedbacks, err := feedbackBase.ListFeedbacks()
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:59:02.440137+06:00")
	if err != nil {
		t.Error(err)
	}

	feedbacktest:=Feedback{
		Id:feedbacks[0].Id,
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"very tasty BESH",
		Created: crtd,
		DealId:1,
		Anon:1,}
	if *feedbacks[0] != feedbacktest{
		fmt.Println(feedbacktest)
		fmt.Println(*feedbacks[0])
		t.Error("Title is not equal")
	}
}

func TestPostgreFeedbackBase_DeleteFeedback(t *testing.T) {
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
	feedbackBase, err = NewPostgreFeedbackBase(db)
	if err != nil {
		t.Error(err)
	}
	crtd, err := time.Parse("2006-01-02T15:04:05Z07:00", "2020-04-12T22:57:57.244503+06:00")
	if err != nil {
		t.Error(err)
	}
	feedback:=&Feedback{
		ConsumerId:5,
		ProducerId:3,
		ConsumerName:"antoshka",
		Value: 2,
		Text:"very tasty BESH",
		Created: crtd,
		DealId:1,
		Anon:1,}
	feedback,err = feedbackBase.CreateFeedback(feedback)
	if err != nil {
		t.Error(err)
	}
	feedbackBase.DeleteFeedback(feedback.Id)
	_, err = feedbackBase.GetFeedback(feedback)
	if err == nil {
		t.Errorf("Not deleted")
	}
}
