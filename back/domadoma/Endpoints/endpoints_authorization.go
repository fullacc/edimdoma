package Endpoints

import (
	"encoding/json"
	"errors"
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/SMS"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthorizationEndpoints interface {
	RegisterUser() func(c *gin.Context)

	LoginUser() func(c *gin.Context)

	LogoutUser() func(c *gin.Context)

	ChangePassword() func(c *gin.Context)

	CheckPhone() func(c *gin.Context)

	CheckCode() func(c *gin.Context)
}

func NewAuthorizationEndpoints(authorizationBase Authorization.AuthorizationBase, smsBase SMS.SMSBase, userBase User.UserBase) AuthorizationEndpoints {
	return &AuthorizationEndpointsFactory{authorizationBase: authorizationBase, smsBase: smsBase, userBase: userBase}
}

type AuthorizationEndpointsFactory struct {
	authorizationBase Authorization.AuthorizationBase
	userBase          User.UserBase
	smsBase           SMS.SMSBase
}

func (f AuthorizationEndpointsFactory) RegisterUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		currtoken, err := f.authorizationBase.GetRegistrationToken(c.Request.Header.Get("Token"))
		if err != nil || currtoken == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if currtoken.Code != "goodtogo" {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		var user Authorization.UserRegister
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": "Provided data is in wrong format"})
			return
		}

		user.UserName = strings.ToLower(user.UserName)
		matched, err := Authorization.Validator(Authorization.Usrnm, user.UserName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't validate username"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db error"})
			return
		}

		matched, err = Authorization.Validator(Authorization.Phn, currtoken.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't validate phone"})
			return
		}

		if !matched {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Phone input"})
			return
		}

		usertocheck = &User.User{Phone: currtoken.Phone}
		usertocheck, err = f.userBase.GetUser(usertocheck)
		if usertocheck != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Such phone exists"})
			return
		}

		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db error"})
			return
		}

		pwd := []byte(user.Password)
		newuser := &User.User{}
		hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't make your password safe"})
			return
		}

		newuser.PasswordHash = hash
		newuser.UserName = user.UserName
		newuser.Role = Authorization.Regular
		newuser.Phone = currtoken.Phone
		newuser.RatingN = 0
		newuser.RatingTotal = 0
		newuser.Rating = 0
		newuser.Name = "Sultan"
		newuser.Surname = "Nur"
		newuser.City = "Almaty"
		result, err := f.userBase.CreateUser(newuser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't create user"})
			return
		}

		token, err := Authorization.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Can't make you safe"})
			return
		}

		input := &Authorization.AuthToken{Permission: newuser.Role, Token: token, UserId: result.Id}
		data, err := json.Marshal(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "I just can't"})
			return
		}

		err = f.authorizationBase.SetToken(input.Token, data, 4*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Can't save token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"User": result, "Token": input.Token})
	}
}

func (f AuthorizationEndpointsFactory) LoginUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user Authorization.UserLogin
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": "Provided data is in wrong format"})
			return
		}
		user.Login = strings.ToLower(user.Login)

		lookupuser := &User.User{}
		matched, err := Authorization.Validator(Authorization.Phn, user.Login)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't validate phone"})
			return
		}
		if matched {
			lookupuser.Phone = user.Login
		} else {
			matched, err = Authorization.Validator(Authorization.Usrnm, user.Login)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't validate username"})
				return
			}
			if matched {
				lookupuser.UserName = user.Login
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid login input"})
				return
			}
		}

		lookupuser, err = f.userBase.GetUser(lookupuser)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such user in system": user.Login})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Debug": err.Error()})
			return
		}

		err = bcrypt.CompareHashAndPassword(lookupuser.PasswordHash, []byte(user.Password))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Wrong password"})
			return
		}

		token, err := Authorization.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Can't make you safe"})
			return
		}

		input := &Authorization.AuthToken{Permission: lookupuser.Role, Token: token, UserId: lookupuser.Id}
		data, err := json.Marshal(input)
		err = f.authorizationBase.SetToken(input.Token, data, 5*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't make you safe"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Token": input.Token})
	}
}

func (f AuthorizationEndpointsFactory) LogoutUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		err = f.authorizationBase.DeleteToken(curruser.Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't delete token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Logged out": curruser.UserId})
	}
}

func (f AuthorizationEndpointsFactory) ChangePassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		id := c.Param("userid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": "No id provided"})
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

		pass := Authorization.UserChangePassword{}
		err = c.ShouldBindJSON(&pass)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided data is in wrong format"})
			return
		}

		user := &User.User{Id: intid}
		user, err = f.userBase.GetUser(user)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such user in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db error"})
			return
		}

		err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pass.OldPassword))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Wrong password"})
			return
		}

		newpwd := []byte(pass.NewPassword)
		user.PasswordHash, err = bcrypt.GenerateFromPassword(newpwd, bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't make your password safe"})
			return
		}

		_, err = f.userBase.UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't update user"})
		}
		c.JSON(http.StatusOK, gin.H{"Changed password for": intid})
	}
}

func (f AuthorizationEndpointsFactory) CheckPhone() func(c *gin.Context) {
	return func(c *gin.Context) {
		number := Authorization.RegistrationPhone{}
		err := c.ShouldBindJSON(&number)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "wrong data"})
			return
		}

		valid, err := Authorization.Validator(Authorization.Phn, number.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal error"})
			return
		}

		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid phone number"})
			return
		}

		user := &User.User{Phone: number.Phone}
		user, err = f.userBase.GetUser(user)
		if user != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Such phone exists"})
			return
		}

		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db error"})
			return
		}

		sms := SMS.SMS{Phone: number.Phone}
		sentSMS, err := f.smsBase.SendSMS(sms)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't send sms"})
			return
		}

		token, err := Authorization.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't make you safe"})
			return
		}

		input := &Authorization.RegistrationToken{Token: token, Phone: sentSMS.Phone, Code: sentSMS.Code}
		data, err := json.Marshal(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "I just can't;("})
			return
		}

		err = f.authorizationBase.SetToken(input.Token, data, 5*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "I just can't save token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Token": input.Token})
	}
}

func (f AuthorizationEndpointsFactory) CheckCode() func(c *gin.Context) {
	return func(c *gin.Context) {
		currtoken, err := f.authorizationBase.GetRegistrationToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		codetocheck := &Authorization.RegistrationCode{}
		err = c.ShouldBindJSON(&codetocheck)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided Code is in wrong format"})
			return
		}

		matched, err := Authorization.Validator(Authorization.Cd, codetocheck.Code)
		if err != nil || !matched {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't validate Code"})
			return
		}

		if codetocheck.Code != currtoken.Code {
			_ = f.authorizationBase.DeleteToken(currtoken.Token)
			c.JSON(http.StatusForbidden, gin.H{"Error": "Wrong Code, removing your token"})
			return
		}

		token, err := Authorization.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't make you safe"})
			return
		}

		input := &Authorization.RegistrationToken{Token: token, Phone: currtoken.Phone, Code: "goodtogo"}
		data, err := json.Marshal(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "I just can't;("})
			return
		}

		err = f.authorizationBase.SetToken(input.Token, data, 5*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "I just can't save token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Token": input.Token})
	}
}
