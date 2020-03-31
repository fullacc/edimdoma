package Feedback

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreFeedbackBase(configfile *domadoma.ConfigFile) (FeedbackBase, error) {

	db := pg.Connect(&pg.Options{
		Database: configfile.PgDbName,
		Addr: configfile.PgDbHost + ":" + configfile.PgDbPort,
		User: configfile.PgDbUser,
		Password: configfile.PgDbPassword,
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
	for _, model := range []interface{}{(*Feedback)(nil)} {
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

func (p *postgreFeedbackBase) ListProducerFeedbacks(id int) ([]*Feedback, error) {
	var feedbacks []*Feedback
	err := p.db.Model(&feedbacks).Where("Producer_Id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (p *postgreFeedbackBase) UpdateFeedback( feedback *Feedback) (*Feedback, error) {
	err := p.db.Update(feedback)
	if err != nil {
		return nil,err
	}
	return feedback, nil
}

func (p *postgreFeedbackBase) DeleteFeedback(id int) error {
	feedback := &Feedback{Id: id}
	err := p.db.Delete(feedback)
	if err != nil {
		return err
	}
	return nil
}