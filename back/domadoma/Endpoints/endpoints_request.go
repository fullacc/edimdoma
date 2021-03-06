package Endpoints

import (
	"errors"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/Request"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"net/http"
	"strconv"
	"time"
)

type RequestEndpoints interface {
	GetRequest() func(c *gin.Context)

	CreateRequest() func(c *gin.Context)

	ListRequests() func(c *gin.Context)

	UpdateRequest() func(c *gin.Context)

	DeleteRequest() func(c *gin.Context)
}

func NewRequestEndpoints(requestBase Request.RequestBase, authorizationBase Authorization.AuthorizationBase, userBase User.UserBase) RequestEndpoints {
	return &RequestEndpointsFactory{requestBase: requestBase, authorizationBase: authorizationBase, userBase: userBase}
}

type RequestEndpointsFactory struct {
	authorizationBase Authorization.AuthorizationBase
	requestBase       Request.RequestBase
	userBase 		  User.UserBase
}

func (f RequestEndpointsFactory) GetRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id := c.Param("requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		request := &Request.Request{Id: intid}
		request, err = f.requestBase.GetRequest(request)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		c.JSON(http.StatusOK, request)
	}
}

func (f RequestEndpointsFactory) CreateRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		id := c.Param("consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		userid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != userid {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		request := Request.Request{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided data format is wrong"})
			return
		}

		user := &User.User{Id:userid}
		user, err  = f.userBase.GetUser(user)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": userid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		request.ConsumerId = user.Id
		request.ConsumerName = user.UserName
		request.Created = time.Now()
		result, err := f.requestBase.CreateRequest(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't create request"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func (f RequestEndpointsFactory) ListRequests() func(c *gin.Context) {
	return func(c *gin.Context) {
		//curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
		//	return
		//}
		//
		//if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
		//	c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
		//	return
		//}
		var err error
		var requests []*Request.Request
		id := c.Param("consumerid")
		if len(id) == 0 {
			requests, err = f.requestBase.ListRequests()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find requests"})
				return
			}
		} else {
			intid, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
				return
			}
			requests, err = f.requestBase.ListConsumerRequests(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find requests"})
				return
			}
		}
		c.JSON(http.StatusOK, requests)
	}
}

func (f RequestEndpointsFactory) UpdateRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		id := c.Param("consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		userid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != userid {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id = c.Param("requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		requesttocheck := &Request.Request{Id: intid}
		requesttocheck, err = f.requestBase.GetRequest(requesttocheck)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		request := &Request.Request{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided data is in wrong format"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && requesttocheck.ConsumerId != userid {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		if request.FoodName == "" {
			request.FoodName = requesttocheck.FoodName
		}

		if request.Description == "" {
			request.Description = requesttocheck.Description
		}

		if request.Type == domadoma.Null {
			request.Type = requesttocheck.Type
		}

		if request.Myaso == domadoma.Null {
			request.Myaso = requesttocheck.Myaso
		}

		if request.Halal == domadoma.Null {
			request.Halal = requesttocheck.Halal
		}

		if request.Vegan == domadoma.Null {
			request.Vegan = requesttocheck.Vegan
		}

		if request.Spicy == domadoma.Null {
			request.Spicy = requesttocheck.Spicy
		}

		if request.Location == nil {
			request.Location = requesttocheck.Location
		}

		if request.Price == 0 {
			request.Price = requesttocheck.Price
		}

		if request.Quantity == 0 {
			request.Quantity = requesttocheck.Quantity
		}

		request.Id = requesttocheck.Id
		request.ConsumerId = requesttocheck.ConsumerId
		request.ConsumerName = requesttocheck.ConsumerName
		request.Created = requesttocheck.Created

		result, err := f.requestBase.UpdateRequest(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't update request"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (f RequestEndpointsFactory) DeleteRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		id := c.Param("consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		userid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != userid {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id = c.Param("requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		request := &Request.Request{Id: intid}
		request, err = f.requestBase.GetRequest(request)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && userid != request.ConsumerId {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		err = f.requestBase.DeleteRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't delete request"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deletedid": intid})
	}
}
