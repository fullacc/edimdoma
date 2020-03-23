package Feedback

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
)

func NewPostgreFeedbackBase(configfile *domadoma.ConfigFile) (FeedbackBase, error) {

	db := pg.Connect(&pg.Options{
		Database: configfile.Name,
		Addr: configfile.DbHost + ":" + configfile.DbPort,
		User: "postgres",
		Password: configfile.Password,
	})

	err := domadoma.createSchema(db)
	if err != nil {
		return nil, err
	}
	return &domadoma.postgreBase{db: db}, nil
}

func (p *domadoma.postgreBase) CreateFeedback(feedback *Feedback) (*Feedback, error) {
	err := p.db.Insert(feedback)
	if err != nil {
		return nil,err
	}
	return feedback,nil
}

func (p *domadoma.postgreBase) GetFeedback(id int) (*Feedback, error) {
	feedback := &Feedback{Id: id}
	err := p.db.Select(&feedback)
	if err != nil {
		return nil, err
	}
	return feedback, nil
}

func (p *domadoma.postgreBase) ListFeedbacks() ([]*Feedback, error) {
	var feedbacks []*Feedback
	err := p.db.Select(feedbacks)
	if err != nil {
		return nil, err
	}
	return feedbacks,nil
}

func (p *domadoma.postgreBase) UpdateFeedback(id int, feedback *Feedback) (*Feedback, error) {
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

func (p *domadoma.postgreBase) DeleteFeedback(id int) error {
	feedback := &Feedback{Id: id}
	err := p.db.Delete(feedback)
	if err != nil {
		return err
	}
	return nil
}