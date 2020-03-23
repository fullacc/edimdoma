package Request

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RequestEndpoints interface{
	GetRequest() func(c *gin.Context)

	CreateRequest() func(c *gin.Context)

	ListRequests() func(c *gin.Context)

	ListConsumerRequests() func(c *gin.Context)

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
		CHECKIFAUTHORIZED
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
		CHECKIFAUTHORIZED
		var request *Request
		if err := c.ShouldBindJSON(&request); err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest,gin.H{"Error: ": err.Error()})
				return
			}
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
		CHECKIFAUTHORIZED
		var requests []*Request
		requests, err := f.requestBase.ListRequests()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,requests)
	}
}

func (f RequestEndpointsFactory) ListConsumerRequests() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKIFAUTHORIZED
		id := c.Param("consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		var requests []*Request
		requests, err = f.requestBase.ListConsumerRequests(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,requests)
	}
}

func (f RequestEndpointsFactory) UpdateRequest() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKAUTHORIZED
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
		request := &Request{}
		if err := c.ShouldBindJSON(request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
			return
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
		CHECKAUTHORIZED
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
		err = f.requestBase.DeleteRequest(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,gin.H{"deleted":intid})
	}
}