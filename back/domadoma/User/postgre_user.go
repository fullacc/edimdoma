package User

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
)

func NewPostgreUserBase(configfile *domadoma.ConfigFile) (UserBase, error) {

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


func (p *domadoma.postgreBase) CreateUser(user *User) (*User, error) {
	err := p.db.Insert(user)
	if err != nil {
		return nil,err
	}
	return user,nil
}

func (p *domadoma.postgreBase) GetUser(user *User) (*User, error) {
	err := p.db.Select(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *domadoma.postgreBase) ListUsers() ([]*User, error) {
	var users []*User
	err := p.db.Select(users)
	if err != nil {
		return nil, err
	}
	return users,nil
}

func (p *domadoma.postgreBase) UpdateUser(id int, user *User) (*User, error) {
	user1 := &User{Id: id}
	err := p.db.Select(user1)
	if err != nil {
		return nil,err
	}
	user1 = user
	err = p.db.Update(user1)
	if err != nil {
		return nil,err
	}
	return user1, nil
}

func (p *domadoma.postgreBase) DeleteUser(id int) error {
	user := &User{Id: id}
	err := p.db.Delete(user)
	if err != nil {
		return err
	}
	return nil
}