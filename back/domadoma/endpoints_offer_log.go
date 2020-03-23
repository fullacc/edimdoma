package domadoma

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OfferLogEndpoints interface{
	GetOfferLog() func(c *gin.Context)

	CreateOfferLog() func(c *gin.Context)

	ListOfferLogs() func(c *gin.Context)

	ListProducerOfferLogs() func(c *gin.Context)

	UpdateOfferLog() func(c *gin.Context)

	DeleteOfferLog() func(c *gin.Context)
}

func NewOfferLogEndpoints(offerLogBase OfferLogBase) OfferLogEndpoints {
	return &OfferLogEndpointsFactory{offerLogBase: offerLogBase}
}

type OfferLogEndpointsFactory struct{
	offerLogBase OfferLogBase
}

func (f OfferLogEndpointsFactory) GetOfferLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKIFAUTHORIZED
		id := c.Param( "offerLogid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":"No id provided"})
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		offerLog, err := f.offerLogBase.GetOfferLog(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,offerLog)
	}
}

func (f OfferLogEndpointsFactory) CreateOfferLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKIFAUTHORIZED
		var offerLog *OfferLog
		if err := c.ShouldBindJSON(&offerLog); err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest,gin.H{"Error: ": err.Error()})
				return
			}
		}
		result, err := f.offerLogBase.CreateOfferLog(offerLog)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,result)
	}
}

func (f OfferLogEndpointsFactory) ListOfferLogs() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKIFAUTHORIZED
		var offerLogs []*OfferLog
		offerLogs, err := f.offerLogBase.ListOfferLogs()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,offerLogs)
	}
}

func (f OfferLogEndpointsFactory) UpdateOfferLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKAUTHORIZED
		id := c.Param("offerLogid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		offerLog := &OfferLog{}
		if err := c.ShouldBindJSON(offerLog); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
			return
		}
		offerLog, err = f.offerLogBase.UpdateOfferLog(intid, offerLog)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,offerLog)
	}
}

func (f OfferLogEndpointsFactory) DeleteOfferLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKAUTHORIZED
		id := c.Param("offerLogid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		err = f.offerLogBase.DeleteOfferLog(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,gin.H{"deleted":intid})
	}
}
