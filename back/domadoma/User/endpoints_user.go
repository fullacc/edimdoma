package User

import (
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

type UserEndpoints interface{
	GetUser() func(c *gin.Context)

	CreateUser() func(c *gin.Context)

	ListUsers() func(c *gin.Context)

	UpdateUser() func(c *gin.Context)

	DeleteUser() func(c *gin.Context)

}

func NewUserEndpoints(userBase UserBase, authorizationBase Authorization.AuthorizationBase) UserEndpoints {
	return &EndpointsFactory{userBase: userBase, authorizationBase:authorizationBase}
}

type EndpointsFactory struct{
	authorizationBase Authorization.AuthorizationBase
	userBase UserBase
}

func (f EndpointsFactory) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager{
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		user := User{}
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": "Provided data is incorrect"})
			return
		}
		user.PasswordHash = "password"
		pwd := []byte(user.PasswordHash)
		hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't make your password safe"})
			return
		}
		user.PasswordHash = hash
		result, err := f.userBase.CreateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't create user"})
			return
		}
		c.JSON(http.StatusCreated,gin.H{"Created":result})
	}
}

func (f EndpointsFactory) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		id := c.Param( "userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		user := &User{Id: intid}
		user, err = f.userBase.GetUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find user"})
			return
		}

		user.PasswordHash=".!."
		c.JSON(http.StatusOK,user)
	}
}

func (f EndpointsFactory) ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager {
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}

		var users []*User
		users, err = f.userBase.ListUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Can't find users"})
			return
		}

		c.JSON(http.StatusCreated,users)
	}
}

func (f EndpointsFactory) UpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != intid{
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		usertocheck := &User{Id:intid}
		usertocheck, err = f.userBase.GetUser(usertocheck)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find user"})
			return
		}

		user := &User{}
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": "Wrong data incoming"})
			return
		}

		user.PasswordHash = usertocheck.PasswordHash
		user.RatingN = usertocheck.RatingN
		user.RatingTotal = usertocheck.RatingTotal
		user.Rating = usertocheck.RatingTotal / usertocheck.RatingN
		user.Id = usertocheck.Id
		if user.UserName == "" {
			user.UserName = usertocheck.UserName
		} else {
			user.UserName = strings.ToLower(user.UserName)
			matched, err := Authorization.Validator(Authorization.Usrnm,user.UserName)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error": "Couldn't Validate"})
				return
			}

			if !matched{
				c.JSON(http.StatusBadRequest,gin.H{"Error":"Invalid Phone input"})
				return
			}
		}
		if user.Name == "" {
			user.Name=usertocheck.Name
		}
		if user.City == "" {
			user.City=usertocheck.City
		}
		if user.Phone == "" {
			user.Phone=usertocheck.Phone
		} else {
			matched, err := Authorization.Validator(Authorization.Phn,user.Phone)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error": "Couldn't Validate"})
				return
			}

			if !matched{
				c.JSON(http.StatusBadRequest,gin.H{"Error":"Invalid Phone input"})
				return
			}
		}
		if user.Surname == "" {
			user.Surname=usertocheck.Surname
		}

		_, err = f.userBase.UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't update user"})
			return
		}
		c.JSON(http.StatusOK,gin.H{"udpated user":intid})
	}
}

func (f EndpointsFactory) DeleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager {
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != intid{
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		err = f.userBase.DeleteUser(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't Delete user"})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deleted":intid})
	}
}

