package OfferLog

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OfferLogEndpoints interface{
	GetOfferLog() func(c *gin.Context)

//	CreateOfferLog() func(c *gin.Context)

	ListOfferLogs() func(c *gin.Context)

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
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.Permission != User.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id := c.Param( "offerlogid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		offer, err := f.offerLogBase.GetOfferLog(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.UserId != offer.ProducerId {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		c.JSON(http.StatusOK,offer)
	}
}
/*
func (f OfferLogEndpointsFactory) CreateOfferLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKIFAUTHORIZED
		var offerLog OfferLog
		if err := c.ShouldBindJSON(&offerLog); err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest,gin.H{"Error ": err.Error()})
				return
			}
		}
		result, err := f.offerLogBase.CreateOfferLog(&offerLog)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,result)
	}
}
*/
func (f OfferLogEndpointsFactory) ListOfferLogs() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.Permission != User.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		var offerLogs []*OfferLog
		idp := c.Param("producerid")
		if (curruser.Permission == User.Admin || curruser.Permission == User.Manager)&&len(idp) == 0 {
			offerLogs, err = f.offerLogBase.ListOfferLogs()
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
				return
			}
		} else{
			if len(idp) != 0 {
				intid, err := strconv.Atoi(idp)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
					return
				}

				offerLogs, err = f.offerLogBase.ListProducerOfferLogs(intid)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
					return
				}

			} else {
				c.JSON(http.StatusForbidden,gin.H{"Error ": "Not allowed"})
				return
			}
		}
		c.JSON(http.StatusOK,offerLogs)
	}
}

func (f OfferLogEndpointsFactory) DeleteOfferLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.Permission != User.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id := c.Param("offerLogid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		offerLogtocheck, err := f.offerLogBase.GetOfferLog(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.UserId != offerLogtocheck.ProducerId{
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		err = f.offerLogBase.DeleteOfferLog(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deletedid":intid})
	}
}
