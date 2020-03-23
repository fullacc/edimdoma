package Feedback

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreFeedbackBase(configfile *domadoma.ConfigFile) (FeedbackBase, error) {

	db := pg.Connect(&pg.Options{
		Database: configfile.Name,
		Addr: configfile.DbHost + ":" + configfile.DbPort,
		User: "postgres",
		Password: configfile.Password,
	})

	err := createSchema(db)
	if err != nil {
		return nil, err
	}
	return &postgreFeedbackBase{db: db}, nil
}

type postgreFeedbackBase struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*FeedbackBase)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *postgreFeedbackBase) CreateFeedback(feedback *Feedback) (*Feedback, error) {
	err := p.db.Insert(feedback)
	if err != nil {
		return nil,err
	}
	return feedback,nil
}

func (p *postgreFeedbackBase) GetFeedback(id int) (*Feedback, error) {
	feedback := &Feedback{Id: id}
	err := p.db.Select(&feedback)
	if err != nil {
		return nil, err
	}
	return feedback, nil
}

func (p *postgreFeedbackBase) ListFeedbacks() ([]*Feedback, error) {
	var feedbacks []*Feedback
	err := p.db.Select(feedbacks)
	if err != nil {
		return nil, err
	}
	return feedbacks,nil
}

func (p *postgreFeedbackBase) ListProducerDeals(id int) ([]*Feedback, error) {
	var feedbacks []*Feedback
	err := p.db.Model(&feedbacks).Where("Producer_Id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (p *postgreFeedbackBase) UpdateFeedback(id int, feedback *Feedback) (*Feedback, error) {
	feedback1 := &Feedback{Id: id}
	err := p.db.Select(feedback1)
	if err != nil {
		return nil,err
	}
	feedback1 = feedback
	err = p.db.Update(feedback1)
	if err != nil {
		return nil,err
	}
	return feedback1, nil
}

func (p *postgreFeedbackBase) DeleteFeedback(id int) error {
	feedback := &Feedback{Id: id}
	err := p.db.Delete(feedback)
	if err != nil {
		return err
	}
	return nil
}