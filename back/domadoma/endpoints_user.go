package domadoma

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
}



func NewUserEndpoints(userBase UserBase) UserEndpoints {
	return &UserEndpointsFactory{userBase: userBase}
}

type UserEndpointsFactory struct{
	userBase UserBase
}

func Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}

func (f UserEndpointsFactory) LoginUser() func (c *gin.Context) {
	return func(c *gin.Context) {
		var user *UserLogin
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":err.Error()})
			return
		}
		lookupuser := &User{UserName:user.UserName}
		lookupuser, err := f.userBase.GetUser(lookupuser)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(lookupuser.PasswordHash), []byte(user.Pwd))
		if err != nil {
			c.JSON(http.StatusForbidden,gin.H{"Error: ": "Wrong password"})
			return
		}
		input := &UserInfo{Permission:Regular,Token:xid.New().String(),UserId:lookupuser.Id}
		data, err := json.Marshal(input)
		redisClient := Connect()
		err = redisClient.Set(input.Token, data, 120 * time.Minute).Err()
		c.JSON(http.StatusOK,gin.H{"Token":input.Token})
	}
}

func (f UserEndpointsFactory) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := getToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}
		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.Permission !=Regular {
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
		user := &User{Id:intid}
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
		pwd := []byte(user.Password)
		user.Password = ".!."
		hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		user.PasswordHash = hash
		result, err := f.userBase.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		input := &UserInfo{Permission:Regular,Token:xid.New().String(),UserId:user.Id}
		data, err := json.Marshal(input)
		redisClient := Connect()
		err = redisClient.Set(input.Token, data, 120 * time.Minute).Err()
		result.PasswordHash=".!."
		c.JSON(http.StatusCreated,gin.H{"User":result,"Token":input.Token})
	}
}


func (f UserEndpointsFactory) ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := getToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}
		if curruser.Permission != Admin && curruser.Permission != Manager {
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
		curruser, err := getToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}
		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.Permission !=Regular {
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
		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.UserId != intid{
			c.JSON(http.StatusForbidden,gin.H{"Error: ":"Not allowed"})
			return
		}
		user := &User{}
		if err := c.ShouldBindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
			return
		}
		user, err = f.userBase.UpdateUser(intid, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,user)
	}
}

func (f UserEndpointsFactory) DeleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := getToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}
		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.Permission !=Regular {
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
		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.UserId != intid{
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

