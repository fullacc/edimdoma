package Request

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
)

func NewPostgreRequestBase(configfile *domadoma.ConfigFile) (RequestBase, error) {

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


func (p *domadoma.postgreBase) CreateRequest(request *Request) (*Request, error) {
	err := p.db.Insert(request)
	if err != nil {
		return nil,err
	}
	return request,nil
}

func (p *domadoma.postgreBase) GetRequest(id int) (*Request, error) {
	request := &Request{Id: id}
	err := p.db.Select(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (p *domadoma.postgreBase) ListRequests() ([]*Request, error) {
	var requests []*Request
	err := p.db.Select(requests)
	if err != nil {
		return nil, err
	}
	return requests,nil
}

func (p *domadoma.postgreBase) ListConsumerRequests(id int) ([]*Request, error) {
	var requests []*Request
	err := p.db.Model(&requests).Where("Consumer_Id = ?",id).Select()
	if err != nil {
		return nil, err
	}
	return requests,nil
}

func (p *domadoma.postgreBase) UpdateRequest(id int, request *Request) (*Request, error) {
	request1 := &Request{Id: id}
	err := p.db.Select(request1)
	if err != nil {
		return nil,err
	}
	request1 = request
	err = p.db.Update(request1)
	if err != nil {
		return nil,err
	}
	return request1, nil
}

func (p *domadoma.postgreBase) DeleteRequest(id int) error {
	request := &Request{Id: id}
	err := p.db.Delete(request)
	if err != nil {
		return err
	}
	return nil
}