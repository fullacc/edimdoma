package Request

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RequestEndpoints interface{
	GetRequest() func(c *gin.Context)

	CreateRequest() func(c *gin.Context)

	ListRequests() func(c *gin.Context)

	UpdateRequest() func(c *gin.Context)

	DeleteRequest() func(c *gin.Context)

}

func NewRequestEndpoints(requestBase RequestBase) RequestEndpoints {
	return &RequestEndpointsFactory{requestBase: requestBase}
}

type RequestEndpointsFactory struct{
	requestBase RequestBase
}

func (f RequestEndpointsFactory) GetRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		id := c.Param( "requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		request, err := f.requestBase.GetRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,request)
	}
}

func (f RequestEndpointsFactory) CreateRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		var request *Request
		if err := c.ShouldBindJSON(&request); err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest,gin.H{"Error: ": err.Error()})
				return
			}
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.UserId != request.ConsumerId {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		result, err := f.requestBase.CreateRequest(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusCreated,result)
	}
}

func (f RequestEndpointsFactory) ListRequests() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		var requests []*Request
		id := c.Param("consumerid")
		if len(id) == 0 {
			requests, err = f.requestBase.ListRequests()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
				return
			}
		} else {
			intid, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
				return
			}
			requests, err = f.requestBase.ListConsumerRequests(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK,requests)	}
}

func (f RequestEndpointsFactory) UpdateRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		id := c.Param("requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		requesttocheck, err := f.requestBase.GetRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		request := &Request{}
		if err := c.ShouldBindJSON(request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && requesttocheck.ConsumerId != curruser.UserId && request.ConsumerId != requesttocheck.ConsumerId {
			c.JSON(http.StatusForbidden,gin.H{"Error :": "Not allowed"})
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

		request, err = f.requestBase.UpdateRequest(intid, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,request)
	}
}

func (f RequestEndpointsFactory) DeleteRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular{
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		id := c.Param("requestid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		requesttocheck, err := f.requestBase.GetRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.UserId != requesttocheck.ConsumerId{
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		err = f.requestBase.DeleteRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deletedid":intid})
	}
}