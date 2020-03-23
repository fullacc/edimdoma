package Offer

import (
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OfferEndpoints interface{
	GetOffer() func(c *gin.Context)

	CreateOffer() func(c *gin.Context)

	ListOffers() func(c *gin.Context)

	UpdateOffer() func(c *gin.Context)

	DeleteOffer() func(c *gin.Context)

}

func NewOfferEndpoints(offerBase OfferBase) OfferEndpoints {
	return &OfferEndpointsFactory{offerBase: offerBase}
}

type OfferEndpointsFactory struct{
	offerBase OfferBase
}
func (f OfferEndpointsFactory) GetOffer() func(c *gin.Context) {
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

		id := c.Param( "offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		offer, err := f.offerBase.GetOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,offer)
	}
}

func (f OfferEndpointsFactory) CreateOffer() func(c *gin.Context) {
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

		var offer *Offer
		if err := c.ShouldBindJSON(&offer); err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest,gin.H{"Error: ": err.Error()})
				return
			}
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && offer.ProducerId != curruser.UserId {
			c.JSON(http.StatusForbidden,gin.H{"Error :": "Not allowed"})
			return
		}

		result, err := f.offerBase.CreateOffer(offer)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,result)
	}
}

func (f OfferEndpointsFactory) ListOffers() func(c *gin.Context) {
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

		var offers []*Offer
		id := c.Param("producerid")
		if len(id) == 0 {
			offers, err = f.offerBase.ListOffers()
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
			offers, err = f.offerBase.ListProducerOffers(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK,offers)
	}
}


func (f OfferEndpointsFactory) UpdateOffer() func(c *gin.Context) {
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

		id := c.Param("offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		offertocheck, err := f.offerBase.GetOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		offer := &Offer{}
		if err := c.ShouldBindJSON(offer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && offertocheck.ProducerId != curruser.UserId && offer.ProducerId != offertocheck.ProducerId {
			c.JSON(http.StatusForbidden,gin.H{"Error :": "Not allowed"})
			return
		}

		offer, err = f.offerBase.UpdateOffer(intid, offer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,offer)
	}
}

func (f OfferEndpointsFactory) DeleteOffer() func(c *gin.Context) {
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

		id := c.Param("offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		offertocheck, err := f.offerBase.GetOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.UserId != offertocheck.ProducerId{
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		err = f.offerBase.DeleteOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deletedid":intid})
	}
}
