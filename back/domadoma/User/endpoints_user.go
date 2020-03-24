package User

import (
	"encoding/json"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type UserEndpoints interface{
	GetUser() func(c *gin.Context)

	CreateUser() func(c *gin.Context)

	ListUsers() func(c *gin.Context)

	UpdateUser() func(c *gin.Context)

	DeleteUser() func(c *gin.Context)

	LoginUser() func(c *gin.Context)

	LogoutUser() func(c *gin.Context)
}



func NewUserEndpoints(userBase UserBase) UserEndpoints {
	return &UserEndpointsFactory{userBase: userBase}
}

type UserEndpointsFactory struct{
	userBase UserBase
}

func (f UserEndpointsFactory) LoginUser() func (c *gin.Context) {
	return func(c *gin.Context) {
		var user *domadoma.UserLogin
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":err.Error()})
			return
		}

		lookupuser := &User{UserName: user.UserName}
		lookupuser, err := f.userBase.GetUser(lookupuser)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(lookupuser.PasswordHash), []byte(user.Password))
		if err != nil {
			c.JSON(http.StatusForbidden,gin.H{"Error: ": "Wrong password"})
			return
		}

		input := &domadoma.UserInfo{Permission: domadoma.Regular,Token:xid.New().String(),UserId:lookupuser.Id}
		data, err := json.Marshal(input)
		redisClient := domadoma.Connect()
		err = redisClient.Set(input.Token, data, 120 * time.Minute).Err()
		c.JSON(http.StatusOK,gin.H{"Token":input.Token})
	}
}

func (f UserEndpointsFactory) LogoutUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		redisClient := domadoma.Connect()
		_,err = redisClient.Del(curruser.Token).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"Logged out":curruser.UserId})
	}
}

func (f UserEndpointsFactory) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden,gin.H{"Error: ":"Not allowed"})
			return
		}

		id := c.Param( "userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		user := &User{Id: intid}
		user, err = f.userBase.GetUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		user.PasswordHash=".!."
		c.JSON(http.StatusOK,user)
	}
}

func (f UserEndpointsFactory) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user *User
		if err := c.ShouldBindJSON(&user); err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest,gin.H{"Error: ": err.Error()})
				return
			}
		}

		usertocheck := &User{Email:user.Email}
		usertocheck,_ = f.userBase.GetUser(usertocheck)
		if usertocheck != nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":"Such email exists"})
			return
		}
		usertocheck = &User{UserName:user.UserName}
		usertocheck,_ = f.userBase.GetUser(usertocheck)
		if usertocheck != nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":"Such username exists"})
			return
		}
		usertocheck = &User{Phone:user.Phone}
		usertocheck,_ = f.userBase.GetUser(usertocheck)
		if usertocheck != nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":"Such phone exists"})
			return
		}

		pwd := []byte(user.Password)
		user.Password = ".!."
		hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err})
			return
		}

		user.PasswordHash = hash
		result, err := f.userBase.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err})
			return
		}

		input := &domadoma.UserInfo{Permission: domadoma.Regular,Token:xid.New().String(),UserId:user.Id}
		data, err := json.Marshal(input)
		redisClient := domadoma.Connect()
		err = redisClient.Set(input.Token, data, 120 * time.Minute).Err()
		result.PasswordHash=".!."
		c.JSON(http.StatusCreated,gin.H{"User":result,"Token":input.Token})
	}
}

func (f UserEndpointsFactory) ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager {
			c.JSON(http.StatusForbidden,gin.H{"Error :":"Not allowed"})
			return
		}

		var users []*User
		users, err = f.userBase.ListUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusCreated,users)
	}
}

func (f UserEndpointsFactory) UpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden,gin.H{"Error: ":"Not allowed"})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.UserId != intid{
			c.JSON(http.StatusForbidden,gin.H{"Error: ":"Not allowed"})
			return
		}

		usertocheck := &User{Id:intid}
		usertocheck,err = f.userBase.GetUser(usertocheck)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		user := &User{}
		if err := c.ShouldBindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
			return
		}

		user.PasswordHash = usertocheck.PasswordHash
		user.RatingN = usertocheck.RatingN
		user.RatingTotal = usertocheck.RatingTotal
		user.Id = usertocheck.Id
		user.Password = usertocheck.Password
		if user.UserName == "" {
			user.UserName=usertocheck.UserName
		}
		if user.Name == "" {
			user.Name=usertocheck.Name
		}
		if user.City == "" {
			user.City=usertocheck.City
		}
		if user.Email == "" {
			user.Email=usertocheck.Email
		}
		if user.Phone == "" {
			user.Phone=usertocheck.Phone
		}
		if user.Surname == "" {
			user.Surname=usertocheck.Surname
		}

		user, err = f.userBase.UpdateUser(intid, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
			return
		}

		user.PasswordHash = ".!."
		user.Password = ".!."
		c.JSON(http.StatusOK,user)
	}
}

func (f UserEndpointsFactory) DeleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden,gin.H{"Error: ":"Not allowed"})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.UserId != intid{
			c.JSON(http.StatusForbidden,gin.H{"Error: ":"Not allowed"})
			return
		}

		err = f.userBase.DeleteUser(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deleted":intid})
	}
}

