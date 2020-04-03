package Endpoints

import (
	"errors"
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/Offer"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"net/http"
	"strconv"
	"time"
)

type OfferEndpoints interface {
	GetOffer() func(c *gin.Context)

	CreateOffer() func(c *gin.Context)

	ListOffers() func(c *gin.Context)

	UpdateOffer() func(c *gin.Context)

	DeleteOffer() func(c *gin.Context)
}

func NewOfferEndpoints(offerBase Offer.OfferBase, authorizationBase Authorization.AuthorizationBase) OfferEndpoints {
	return &OfferEndpointsFactory{offerBase: offerBase, authorizationBase: authorizationBase}
}

type OfferEndpointsFactory struct {
	authorizationBase Authorization.AuthorizationBase
	offerBase         Offer.OfferBase
}

func (f OfferEndpointsFactory) GetOffer() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}
		offer := &Offer.Offer{Id: intid}
		offer, err = f.offerBase.GetOffer(offer)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		c.JSON(http.StatusOK, offer)
	}
}

func (f OfferEndpointsFactory) CreateOffer() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		id := c.Param("producerid")
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

		offer := Offer.Offer{}
		err = c.ShouldBindJSON(&offer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided data is in wrong format"})
			return
		}

		offer.ProducerId = userid
		offer.Created = time.Now()
		offer.AvailableQuantity = offer.InitialQuantity

		result, err := f.offerBase.CreateOffer(&offer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't create offer"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

func (f OfferEndpointsFactory) ListOffers() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		var offers []*Offer.Offer
		id := c.Param("producerid")
		if len(id) == 0 {
			offers, err = f.offerBase.ListOffers()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find offers"})
				return
			}
		} else {
			intid, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
				return
			}
			offers, err = f.offerBase.ListProducerOffers(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find offers"})
				return
			}
		}
		c.JSON(http.StatusOK, offers)
	}
}

func (f OfferEndpointsFactory) UpdateOffer() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		offertocheck := &Offer.Offer{Id: intid}
		offertocheck, err = f.offerBase.GetOffer(offertocheck)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		offer := &Offer.Offer{}
		err = c.ShouldBindJSON(&offer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided data is in wrong format"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && offertocheck.ProducerId != curruser.UserId && offer.ProducerId != offertocheck.ProducerId {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		if offer.Name == "" {
			offer.Name = offertocheck.Name
		}

		offer.ProducerId = offertocheck.ProducerId

		if offer.Created.IsZero() {
			offer.Created = offertocheck.Created
		}

		if offer.Location == nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't update offer"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (f OfferEndpointsFactory) DeleteOffer() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("offerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		offertocheck := &Offer.Offer{Id: intid}
		offertocheck, err = f.offerBase.GetOffer(offertocheck)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != offertocheck.ProducerId {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		err = f.offerBase.DeleteOffer(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't delete offer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"deletedid": intid})
	}
}
