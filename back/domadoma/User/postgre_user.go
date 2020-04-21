package User

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func NewPostgreUserBase(db *pg.DB) (UserBase, error) {
	for _, model := range []interface{}{(*User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return nil, err
		}
	}
	return &postgreUserBase{db: db}, nil
}

type postgreUserBase struct {
	db *pg.DB
}

func (p *postgreUserBase) CreateUser(user *User) (*User, error) {
	err := p.db.Insert(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *postgreUserBase) GetUser(user *User) (*User, error) {
	err := error(nil)
	if user.Id != 0 {
		err = p.db.Select(user)
	} else {
		if user.UserName != "" {
			err = p.db.Model(user).Where("user_name = ?", user.UserName).Limit(1).Select()
		} else {
			if user.Phone != "" {
				err = p.db.Model(user).Where("phone = ?", user.Phone).Limit(1).Select()
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
	err := p.db.Model(&users).Select()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (p *postgreUserBase) UpdateUser(user *User) (*User, error) {
	err := p.db.Update(user)
	if err != nil {
		return nil, err
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
