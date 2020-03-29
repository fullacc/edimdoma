package Request

import (
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type RequestEndpoints interface{
	GetRequest() func(c *gin.Context)

	CreateRequest() func(c *gin.Context)

	ListRequests() func(c *gin.Context)

	UpdateRequest() func(c *gin.Context)

	DeleteRequest() func(c *gin.Context)

}

func NewRequestEndpoints(requestBase RequestBase, authorizationBase Authorization.AuthorizationBase) RequestEndpoints {
	return &EndpointsFactory{requestBase: requestBase, authorizationBase:authorizationBase}
}

type EndpointsFactory struct{
	authorizationBase Authorization.AuthorizationBase
	requestBase RequestBase
}

func (f EndpointsFactory) GetRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id := c.Param( "requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		request, err := f.requestBase.GetRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find request"})
			return
		}

		c.JSON(http.StatusOK,request)
	}
}

func (f EndpointsFactory) CreateRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		id := c.Param("consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		userid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != userid {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		request := Request{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "Provided data format is wrong"})
			return
		}

		request.ConsumerId = userid
		request.Created = time.Now()
		result, err := f.requestBase.CreateRequest(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't create request"})
			return
		}

		c.JSON(http.StatusCreated,result)
	}
}

func (f EndpointsFactory) ListRequests() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		var requests []*Request
		id := c.Param("consumerid")
		if len(id) == 0 {
			requests, err = f.requestBase.ListRequests()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't find requests"})
				return
			}
		} else {
			intid, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Provided id is not integer"})
				return
			}
			requests, err = f.requestBase.ListConsumerRequests(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error ":"Couldn't find requests"})
				return
			}
		}
		c.JSON(http.StatusOK,requests)	}
}

func (f EndpointsFactory) UpdateRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		id := c.Param("consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		userid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != userid{
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id = c.Param("requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		requesttocheck, err := f.requestBase.GetRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find request"})
			return
		}

		request := &Request{}
		err = c.ShouldBindJSON(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": "Provided data is in wrong format"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && requesttocheck.ConsumerId != userid {
			c.JSON(http.StatusForbidden,gin.H{"Error": "Not allowed"})
			return
		}

		request.Id = requesttocheck.Id

		if request.Food == "" {
			request.Food = requesttocheck.Food
		}

		request.ConsumerId = requesttocheck.ConsumerId

		if request.Created.IsZero() {
			request.Created = requesttocheck.Created
		}

		if request.Location == nil{
			request.Location = requesttocheck.Location
		}

		if request.Price == 0 {
			request.Price = requesttocheck.Price
		}

		if request.Quantity == 0 {
			request.Quantity = requesttocheck.Quantity
		}

		result, err := f.requestBase.UpdateRequest(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't update request"})
			return
		}

		c.JSON(http.StatusOK,result)
	}
}

func (f EndpointsFactory) DeleteRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		id := c.Param("consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		userid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != userid {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id = c.Param("requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		requesttocheck, err := f.requestBase.GetRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find request"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && userid != requesttocheck.ConsumerId{
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		err = f.requestBase.DeleteRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't delete request"})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deletedid":intid})
	}
}