package Endpoints

import (
	"errors"
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/Offer"
	"github.com/fullacc/edimdoma/back/domadoma/Rabbit"
	"github.com/fullacc/edimdoma/back/domadoma/User"
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

func NewOfferEndpoints(offerBase Offer.OfferBase, authorizationBase Authorization.AuthorizationBase, userBase User.UserBase, rabbitBase Rabbit.RabbitBase) OfferEndpoints {
	return &OfferEndpointsFactory{offerBase: offerBase, authorizationBase: authorizationBase, userBase: userBase, rabbitBase: rabbitBase}
}

type OfferEndpointsFactory struct {
	authorizationBase Authorization.AuthorizationBase
	offerBase         Offer.OfferBase
	userBase		  User.UserBase
	rabbitBase Rabbit.RabbitBase
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
		offer.ProducerId = user.Id
		offer.ProducerName = user.UserName
		offer.ProducerRating = user.Rating
		offer.Created = time.Now()
		offer.AvailableQuantity = offer.InitialQuantity

		result, err := f.offerBase.CreateOffer(&offer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't create offer"})
			return
		}
		_, err = f.rabbitBase.CreateRabbit(result.Id)

		c.JSON(http.StatusCreated, result)
	}
}

func (f OfferEndpointsFactory) ListOffers() func(c *gin.Context) {
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
		var err = error(nil)

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

		if offer.FoodName == "" {
			offer.FoodName = offertocheck.FoodName
		}

		if offer.Description == "" {
			offer.Description = offertocheck.Description
		}

		if offer.Type == 0 {
			offer.Type = offertocheck.Type
		}

		if offer.Myaso == 0 {
			offer.Myaso = offertocheck.Myaso
		}

		if offer.Halal == 0 {
			offer.Halal = offertocheck.Halal
		}

		if offer.Vegan == 0 {
			offer.Vegan = offertocheck.Vegan
		}

		if offer.Spicy == 0 {
			offer.Spicy = offertocheck.Spicy
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
		offer.ProducerId = offertocheck.ProducerId
		offer.ProducerName = offertocheck.ProducerName
		offer.ProducerRating = offertocheck.ProducerRating
	 	offer.Created = offertocheck.Created

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
