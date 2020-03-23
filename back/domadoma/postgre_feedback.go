package domadoma

import "github.com/go-pg/pg"

func NewPostgreFeedbackBase(configfile *ConfigFile) (FeedbackBase, error) {

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
	return &postgreBase{db: db}, nil
}

func (p *postgreBase) CreateFeedback(feedback *Feedback) (*Feedback, error) {
	err := p.db.Insert(feedback)
	if err != nil {
		return nil,err
	}
	return feedback,nil
}

func (p *postgreBase) GetFeedback(id int) (*Feedback, error) {
	feedback := &Feedback{Id:id}
	err := p.db.Select(&feedback)
	if err != nil {
		return nil, err
	}
	return feedback, nil
}

func (p *postgreBase) ListFeedbacks() ([]*Feedback, error) {
	var feedbacks []*Feedback
	err := p.db.Select(feedbacks)
	if err != nil {
		return nil, err
	}
	return feedbacks,nil
}

func (p *postgreBase) UpdateFeedback(id int, feedback *Feedback) (*Feedback, error) {
	feedback1 := &Feedback{Id:id}
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

func (p *postgreBase) DeleteFeedback(id int) error {
	feedback := &Feedback{Id: id}
	err := p.db.Delete(feedback)
	if err != nil {
		return err
	}
	return nil
}