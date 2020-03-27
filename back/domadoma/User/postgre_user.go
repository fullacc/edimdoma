package User

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreUserBase(configfile *domadoma.ConfigFile) (UserBase, error) {

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
	err := error(nil)
	if user.Id != 0 {
		err = p.db.Select(&user)
	} else {
		if user.UserName != "" {
			err = p.db.Model(&user).Where("user.user_name = ?",user.UserName).Limit(1).Select()
		} else {
			if user.Phone != "" {
				err = p.db.Model(&user).Where("user.phone = ?",user.Phone).Limit(1).Select()
			} else {
				if user.Email != "" {
					err = p.db.Model(&user).Where("user.email = ?",user.Email).Limit(1).Select()
				}
			}
		}
	}
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

func (p *postgreUserBase) UpdateUser( user *User) (*User, error) {
	err := p.db.Update(user)
	if err != nil {
		return nil,err
	}
	return user, nil
}

func (p *postgreUserBase) DeleteUser(id int) error {
	user := &User{Id: id}
	err := p.db.Delete(user)
	if err != nil {
		return err
	}
	return nil
}