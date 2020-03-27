package User

import (
	"encoding/json"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
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
	
	ChangePassword() func (c *gin.Context)
}



func NewUserEndpoints(userBase UserBase) UserEndpoints {
	return &UserEndpointsFactory{userBase: userBase}
}

type UserEndpointsFactory struct{
	userBase UserBase
}

func (f UserEndpointsFactory) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user UserRegister
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": err.Error()})
			return
		}

		user.Email = strings.ToLower(user.Email)
		matched, err := domadoma.Validator(Eml,user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
			return
		}

		if !matched{
			c.JSON(http.StatusBadRequest,gin.H{"Error":"Invalid Email input"})
			return
		}

		usertocheck := &User{Email:user.Email}
		usertocheck,_ = f.userBase.GetUser(usertocheck)
		if usertocheck != nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"Such email exists"})
			return
		}

		user.UserName = strings.ToLower(user.UserName)
		matched, err = domadoma.Validator(Usrnm,user.UserName)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
			return
		}

		if !matched{
			c.JSON(http.StatusBadRequest,gin.H{"Error":"Invalid username input"})
			return
		}

		usertocheck = &User{UserName:user.UserName}
		usertocheck,_ = f.userBase.GetUser(usertocheck)
		if usertocheck != nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"Such username exists"})
			return
		}

		matched, err = domadoma.Validator(Phn,user.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
			return
		}

		if !matched{
			c.JSON(http.StatusBadRequest,gin.H{"Error":"Invalid Phone input"})
			return
		}

		usertocheck = &User{Phone:user.Phone}
		usertocheck,_ = f.userBase.GetUser(usertocheck)
		if usertocheck != nil{
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"Such phone exists"})
			return
		}

		pwd := []byte(user.Password)
		newuser := &User{}
		hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		newuser.PasswordHash = hash
		newuser.UserName = user.UserName
		newuser.Email = user.Email
		newuser.Phone = user.Phone

		newuser.Name = "Sultan"
		newuser.Surname = "Nur"
		newuser.City = "Almaty"
		result, err := f.userBase.CreateUser(newuser)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		input := &UserInfo{Permission: Regular,Token:xid.New().String(),UserId:result.Id}
		data, err := json.Marshal(input)
		redisClient := domadoma.Connect()
		err = redisClient.Set(input.Token, data, 120 * time.Minute).Err()
		result.PasswordHash=".!."
		c.JSON(http.StatusCreated,gin.H{"User":result,"Token":input.Token})
	}
}

func (f UserEndpointsFactory) GetUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.Permission != Regular {
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
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		user := &User{Id: intid}
		user, err = f.userBase.GetUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		user.PasswordHash=".!."
		c.JSON(http.StatusOK,user)
	}
}

func (f UserEndpointsFactory) ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != Admin && curruser.Permission != Manager {
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}

		var users []*User
		users, err = f.userBase.ListUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		c.JSON(http.StatusCreated,users)
	}
}

func (f UserEndpointsFactory) UpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.Permission != Regular {
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
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.UserId != intid{
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		usertocheck := &User{Id:intid}
		usertocheck,err = f.userBase.GetUser(usertocheck)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		user := User{}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": err.Error()})
			return
		}

		user.PasswordHash = usertocheck.PasswordHash
		user.RatingN = usertocheck.RatingN
		user.RatingTotal = usertocheck.RatingTotal
		user.Id = usertocheck.Id
		if user.UserName == "" {
			user.UserName = usertocheck.UserName
		} else {
			user.UserName = strings.ToLower(user.UserName)
			matched, err := domadoma.Validator(Usrnm,user.UserName)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
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
		if user.Email == "" {
			user.Email=usertocheck.Email
		} else {
			user.Email = strings.ToLower(user.Email)
			matched, err := domadoma.Validator(Eml,user.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
				return
			}

			if !matched{
				c.JSON(http.StatusBadRequest,gin.H{"Error":"Invalid Phone input"})
				return
			}
		}
		if user.Phone == "" {
			user.Phone=usertocheck.Phone
		} else {
			matched, err := domadoma.Validator(Phn,user.Phone)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
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

		result, err := f.userBase.UpdateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func (f UserEndpointsFactory) DeleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.Permission != Regular {
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
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.UserId != intid{
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		err = f.userBase.DeleteUser(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deleted":intid})
	}
}

func (f UserEndpointsFactory) LoginUser() func (c *gin.Context) {
	return func(c *gin.Context) {
		var user UserLogin
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error ":err.Error()})
			return
		}
		user.Login = strings.ToLower(user.Login)

		lookupuser := &User{}
		matched, err := domadoma.Validator(Phn,user.Login)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
			return
		}
		if matched{
			lookupuser.Phone = user.Login
		}else {
			matched, err = domadoma.Validator(Usrnm,user.Login)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
				return
			}
			if matched {
				lookupuser.Email = user.Login
			} else {
				matched, err = domadoma.Validator(Eml,user.Login)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
					return
				}
				if matched {
					lookupuser.UserName = user.Login
				} else {
					c.JSON(http.StatusBadRequest,gin.H{"Error":"Invalid login input"})
					return
				}
			}
		}


		lookupuser, err = f.userBase.GetUser(lookupuser)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(lookupuser.PasswordHash), []byte(user.Password))
		if err != nil {
			c.JSON(http.StatusForbidden,gin.H{"Error ": "Wrong password"})
			return
		}

		input := &UserInfo{Permission: Regular,Token:xid.New().String(),UserId:lookupuser.Id}
		data, err := json.Marshal(input)
		redisClient := domadoma.Connect()
		err = redisClient.Set(input.Token, data, 3 * time.Hour).Err()
		c.JSON(http.StatusOK,gin.H{"Token":input.Token})
	}
}

func (f UserEndpointsFactory) LogoutUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		redisClient := domadoma.Connect()
		_,err = redisClient.Del(curruser.Token).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"Logged out":curruser.UserId})
	}
}

func (f UserEndpointsFactory) ChangePassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.UserId != intid{
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		pass := UserChangePassword{}
		err = c.ShouldBindJSON(&pass)
		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error": err.Error()})
			return
		}

		usertocheck := User{Id:intid}
		user, err := f.userBase.GetUser(&usertocheck)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pass.OldPassword))
		if err != nil {
			c.JSON(http.StatusForbidden,gin.H{"Error ": "Wrong password"})
			return
		}

		newpwd := []byte(pass.NewPassword)
		user.PasswordHash, err = bcrypt.GenerateFromPassword(newpwd, bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		_,err = f.userBase.UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
		}
		c.JSON(http.StatusOK,gin.H{"changed for":intid})
	}
}