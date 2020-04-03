package Endpoints

import (
	"errors"
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

type UserEndpoints interface {
	GetUser() func(c *gin.Context)

	CreateUser() func(c *gin.Context)

	ListUsers() func(c *gin.Context)

	UpdateUser() func(c *gin.Context)

	DeleteUser() func(c *gin.Context)
}

func NewUserEndpoints(userBase User.UserBase, authorizationBase Authorization.AuthorizationBase) UserEndpoints {
	return &UserEndpointsFactory{userBase: userBase, authorizationBase: authorizationBase}
}

type UserEndpointsFactory struct {
	authorizationBase Authorization.AuthorizationBase
	userBase          User.UserBase
}

func (f UserEndpointsFactory) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		user := User.User{}
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided data is incorrect"})
			return
		}

		user.UserName = strings.ToLower(user.UserName)
		matched, err := Authorization.Validator(Authorization.Usrnm, user.UserName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't Validate username"})
			return
		}

		if !matched {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid username input"})
			return
		}

		usertocheck := &User.User{UserName: user.UserName}
		usertocheck, err = f.userBase.GetUser(usertocheck)
		if usertocheck != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Such username exists"})
			return
		}

		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		matched, err = Authorization.Validator(Authorization.Phn, user.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't Validate phone"})
			return
		}

		if !matched {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid phone input"})
			return
		}

		usertocheck = &User.User{Phone: user.Phone}
		usertocheck, err = f.userBase.GetUser(usertocheck)
		if usertocheck != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Such phone exists"})
			return
		}

		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		pwd := []byte("password")
		hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't make your password safe"})
			return
		}
		if c.Query("role") == "manager" {
			user.Role = Authorization.Manager
		} else {
			user.Role = Authorization.Regular
		}
		user.PasswordHash = hash
		user.Rating = 0
		result, err := f.userBase.CreateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't Create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"Created": result})
	}
}

func (f UserEndpointsFactory) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		user := &User.User{Id: intid}
		user, err = f.userBase.GetUser(user)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (f UserEndpointsFactory) ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		var users []*User.User
		users, err = f.userBase.ListUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func (f UserEndpointsFactory) UpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != intid {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		usertocheck := &User.User{Id: intid}
		usertocheck, err = f.userBase.GetUser(usertocheck)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		user := &User.User{}
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong data incoming"})
			return
		}

		user.PasswordHash = usertocheck.PasswordHash
		user.Role = usertocheck.Role
		user.RatingN = usertocheck.RatingN
		user.RatingTotal = usertocheck.RatingTotal
		user.Rating = usertocheck.RatingTotal / usertocheck.RatingN
		user.Id = usertocheck.Id
		if user.UserName == "" {
			user.UserName = usertocheck.UserName
		} else {
			user.UserName = strings.ToLower(user.UserName)
			matched, err := Authorization.Validator(Authorization.Usrnm, user.UserName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't Validate username"})
				return
			}

			if !matched {
				c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Phone input"})
				return
			}

			usertocheckname := &User.User{UserName: user.UserName}
			usertocheckname, err = f.userBase.GetUser(usertocheckname)
			if usertocheckname != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Error": "Such username exists"})
				return
			}

			if err != nil && !errors.Is(err, pg.ErrNoRows) {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
				return
			}

		}
		if user.Name == "" {
			user.Name = usertocheck.Name
		}
		if user.City == "" {
			user.City = usertocheck.City
		}
		if user.Phone == "" {
			user.Phone = usertocheck.Phone
		} else {
			matched, err := Authorization.Validator(Authorization.Phn, user.Phone)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't Validate phone"})
				return
			}

			if !matched {
				c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Phone input"})
				return
			}

			usertocheckphone := &User.User{Phone: user.Phone}
			usertocheckphone, err = f.userBase.GetUser(usertocheckphone)
			if usertocheckphone != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Error": "Such phone exists"})
				return
			}

			if err != nil && !errors.Is(err, pg.ErrNoRows) {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
				return
			}
		}
		if user.Surname == "" {
			user.Surname = usertocheck.Surname
		}

		_, err = f.userBase.UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't Update user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"udpated user": intid})
	}
}

func (f UserEndpointsFactory) DeleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != intid {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		err = f.userBase.DeleteUser(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't Delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deleted": intid})
	}
}
