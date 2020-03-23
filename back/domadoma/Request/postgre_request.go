package Request

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewpostgreRequestBase(configfile *domadoma.ConfigFile) (RequestBase, error) {

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
	return &postgreRequestBase{db: db}, nil
}

type postgreRequestBase struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Request)(nil)} {
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

func (p *postgreRequestBase) CreateRequest(request *Request) (*Request, error) {
	err := p.db.Insert(request)
	if err != nil {
		return nil,err
	}
	return request,nil
}

func (p *postgreRequestBase) GetRequest(id int) (*Request, error) {
	request := &Request{Id: id}
	err := p.db.Select(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (p *postgreRequestBase) ListRequests() ([]*Request, error) {
	var requests []*Request
	err := p.db.Select(requests)
	if err != nil {
		return nil, err
	}
	return requests,nil
}

func (p *postgreRequestBase) ListConsumerRequests(id int) ([]*Request, error) {
	var requests []*Request
	err := p.db.Model(&requests).Where("Consumer_Id = ?",id).Select()
	if err != nil {
		return nil, err
	}
	return requests,nil
}

func (p *postgreRequestBase) UpdateRequest(id int, request *Request) (*Request, error) {
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

func (p *postgreRequestBase) DeleteRequest(id int) error {
	request := &Request{Id: id}
	err := p.db.Delete(request)
	if err != nil {
		return err
	}
	return nil
}