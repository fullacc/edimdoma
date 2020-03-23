package User

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewpostgreUserBase(configfile *domadoma.ConfigFile) (UserBase, error) {

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
	return &postgreUserBase{db: db}, nil
}

type postgreUserBase struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil)} {
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


func (p *postgreUserBase) CreateUser(user *User) (*User, error) {
	err := p.db.Insert(user)
	if err != nil {
		return nil,err
	}
	return user,nil
}

func (p *postgreUserBase) GetUser(user *User) (*User, error) {
	err := p.db.Select(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *postgreUserBase) ListUsers() ([]*User, error) {
	var users []*User
	err := p.db.Select(users)
	if err != nil {
		return nil, err
	}
	return users,nil
}

func (p *postgreUserBase) UpdateUser(id int, user *User) (*User, error) {
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

func (p *postgreUserBase) DeleteUser(id int) error {
	user := &User{Id: id}
	err := p.db.Delete(user)
	if err != nil {
		return err
	}
	return nil
}