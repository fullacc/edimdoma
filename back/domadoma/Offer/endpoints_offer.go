package Offer

import (
	"../Authorization"
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

func NewOfferEndpoints(offerBase OfferBase, authorizationBase Authorization.AuthorizationBase) OfferEndpoints {
	return &EndpointsFactory{offerBase: offerBase, authorizationBase:authorizationBase}
}

type EndpointsFactory struct{
	authorizationBase Authorization.AuthorizationBase
	offerBase OfferBase
}
func (f EndpointsFactory) GetOffer() func(c *gin.Context) {
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

		id := c.Param( "offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		offer, err := f.offerBase.GetOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find offer"})
			return
		}

		c.JSON(http.StatusOK,offer)
	}
}

func (f EndpointsFactory) CreateOffer() func(c *gin.Context) {
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

		offer := Offer{}
		err = c.ShouldBindJSON(&offer)
		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "Provided data is in wrong format"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && offer.ProducerId != curruser.UserId {
			c.JSON(http.StatusForbidden,gin.H{"Error": "Not allowed"})
			return
		}

		result, err := f.offerBase.CreateOffer(&offer)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't create offer"})
			return
		}
		c.JSON(http.StatusCreated,result)
	}
}

func (f EndpointsFactory) ListOffers() func(c *gin.Context) {
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

		var offers []*Offer
		id := c.Param("producerid")
		if len(id) == 0 {
			offers, err = f.offerBase.ListOffers()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't find offers"})
				return
			}
		} else {
			intid, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Provided id is not integer"})
				return
			}
			offers, err = f.offerBase.ListProducerOffers(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't find offers"})
				return
			}
		}
		c.JSON(http.StatusOK,offers)
	}
}


func (f EndpointsFactory) UpdateOffer() func(c *gin.Context) {
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

		id := c.Param("offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		offertocheck, err := f.offerBase.GetOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find offer"})
			return
		}

		offer := &Offer{}
		err = c.ShouldBindJSON(&offer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": "Provided data is in wrong format"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && offertocheck.ProducerId != curruser.UserId && offer.ProducerId != offertocheck.ProducerId {
			c.JSON(http.StatusForbidden,gin.H{"Error": "Not allowed"})
			return
		}

		if offer.Food == "" {
			offer.Food = offertocheck.Food
		}

		offer.ProducerId = offertocheck.ProducerId

		if offer.Created.IsZero() {
			offer.Created = offertocheck.Created
		}

		if offer.Location == nil{
			offer.Location = offertocheck.Location
		}

		if offer.Price == 0 {
			offer.Price = offertocheck.Price
		}

		if offer.InitialQuantity == 0 {
			offer.InitialQuantity = offertocheck.InitialQuantity
		}

		if offer.AvailableQuantity == 0 {
			offer.AvailableQuantity = offertocheck.AvailableQuantity
		}

		offer.Id = offertocheck.Id

		result, err := f.offerBase.UpdateOffer(offer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't update offer"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (f EndpointsFactory) DeleteOffer() func(c *gin.Context) {
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

		id := c.Param("offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		offertocheck, err := f.offerBase.GetOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find offer"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != offertocheck.ProducerId{
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		err = f.offerBase.DeleteOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't delete offer"})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deletedid":intid})
	}
}
