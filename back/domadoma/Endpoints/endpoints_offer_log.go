package Endpoints

import (
	"errors"
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/OfferLog"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"net/http"
	"strconv"
)

type OfferLogEndpoints interface{
	GetOfferLog() func(c *gin.Context)

//	CreateOfferLog() func(c *gin.Context)

	ListOfferLogs() func(c *gin.Context)

	DeleteOfferLog() func(c *gin.Context)
}

func NewOfferLogEndpoints(offerLogBase OfferLog.OfferLogBase, authorizationBase Authorization.AuthorizationBase) OfferLogEndpoints {
	return &OfferLogEndpointsFactory{offerLogBase: offerLogBase, authorizationBase:authorizationBase}
}

type OfferLogEndpointsFactory struct{
	authorizationBase Authorization.AuthorizationBase
	offerLogBase      OfferLog.OfferLogBase
}

func (f OfferLogEndpointsFactory) GetOfferLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param( "offerlogid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": "Provided id is not integer"})
			return
		}

		offer := &OfferLog.OfferLog{Id:intid}
		offer, err = f.offerLogBase.GetOfferLog(offer)
		if err != nil && errors.Is(err,pg.ErrNoRows){
			c.JSON(http.StatusNotFound,gin.H{"No such id in system":intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": "Db Error"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != offer.ProducerId {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
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
		err := c.ShouldBindJSON(&offerLog)
		if err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest,gin.H{"Error": err.Error()})
				return
			}
		}
		result, err := f.offerLogBase.CreateOfferLog(&offerLog)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,result)
	}
}
*/
func (f OfferLogEndpointsFactory) ListOfferLogs() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		var offerLogs []*OfferLog.OfferLog
		idp := c.Param("producerid")
		if (curruser.Permission == Authorization.Admin || curruser.Permission == Authorization.Manager)&&len(idp) == 0 {
			offerLogs, err = f.offerLogBase.ListOfferLogs()
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error": "Couldn't find Offer Logs"})
				return
			}
		} else{
			if len(idp) != 0 {
				intid, err := strconv.Atoi(idp)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error": "Provided id is not integer"})
					return
				}

				offerLogs, err = f.offerLogBase.ListProducerOfferLogs(intid)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error": "Couldn't find offer logs"})
					return
				}

			} else {
				c.JSON(http.StatusForbidden,gin.H{"Error": "Not allowed"})
				return
			}
		}
		c.JSON(http.StatusOK,offerLogs)
	}
}

func (f OfferLogEndpointsFactory) DeleteOfferLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("offerlogid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": "Provided id is not integer"})
			return
		}

		offerLogtocheck := &OfferLog.OfferLog{Id:intid}
		offerLogtocheck, err = f.offerLogBase.GetOfferLog(offerLogtocheck)
		if err != nil && errors.Is(err,pg.ErrNoRows){
			c.JSON(http.StatusNotFound,gin.H{"No such id in system":intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Db Error"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != offerLogtocheck.ProducerId{
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		err = f.offerLogBase.DeleteOfferLog(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error": "Couldn't delete offer log"})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deletedid":intid})
	}
}
