package Endpoints

import (
	"errors"
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/Deal"
	"github.com/fullacc/edimdoma/back/domadoma/Offer"
	"github.com/fullacc/edimdoma/back/domadoma/OfferLog"
	"github.com/fullacc/edimdoma/back/domadoma/Request"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"net/http"
	"strconv"
	"time"
)

type DealEndpoints interface {
	CreateDeal() func(c *gin.Context)

	GetDeal() func(c *gin.Context)

	ListDeals() func(c *gin.Context)

	UpdateDeal() func(c *gin.Context)

	DeleteDeal() func(c *gin.Context)

	CompleteDeal() func(c *gin.Context)
}

func NewDealEndpoints(dealBase Deal.DealBase, authorizationBase Authorization.AuthorizationBase, offerBase Offer.OfferBase, offerLogBase OfferLog.OfferLogBase, requestBase Request.RequestBase, userBase User.UserBase) DealEndpoints {
	return &DealEndpointsFactory{dealBase: dealBase, authorizationBase: authorizationBase, offerBase: offerBase, offerLogBase: offerLogBase, requestBase: requestBase, userBase: userBase}
}

type DealEndpointsFactory struct {
	authorizationBase Authorization.AuthorizationBase
	dealBase          Deal.DealBase
	offerBase         Offer.OfferBase
	offerLogBase      OfferLog.OfferLogBase
	requestBase       Request.RequestBase
	userBase    User.UserBase
}

func (f DealEndpointsFactory) GetDeal() func(c *gin.Context) {
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

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		deal := &Deal.Deal{Id: intid}
		deal, err = f.dealBase.GetDeal(deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find deal"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != deal.ConsumerId && curruser.UserId != deal.ProducerId {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		c.JSON(http.StatusOK, deal)
	}
}

func (f DealEndpointsFactory) CreateDeal() func(c *gin.Context) {
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

		reqid := c.Param("requestid")
		offid := c.Param("offerid")

		deal := Deal.Deal{}
		err = c.ShouldBindJSON(&deal)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided data is in wrong format"})
			return
		}

		if deal.Price < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong price"})
			return
		}

		if len(reqid) != 0 {
			intid, err := strconv.Atoi(reqid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
				return
			}

			request := &Request.Request{Id: intid}
			request, err = f.requestBase.GetRequest(request)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find request"})
				return
			}

			id := c.Param("producerid")
			producerid, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
				return
			}

			if (curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != producerid) || producerid == request.ConsumerId {
				c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
				return
			}

			user := &User.User{Id:producerid}
			user, err  = f.userBase.GetUser(user)
			if err != nil && errors.Is(err, pg.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"No such id in system": producerid})
				return
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
				return
			}

			deal.Quantity = request.Quantity
			deal.FoodName = request.FoodName
			deal.ConsumerId = request.ConsumerId
			deal.ConsumerName = request.ConsumerName
			deal.ProducerId = user.Id
			deal.ProducerName = user.UserName
			deal.Type = request.Type
			deal.Myaso = request.Myaso
			deal.Halal = request.Halal
			deal.Vegan = request.Vegan
			deal.Spicy = request.Spicy
			deal.Description = request.Description


			err = f.requestBase.DeleteRequest(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't delete request"})
				return
			}

		} else {
			if len(offid) != 0 {
				intid, err := strconv.Atoi(offid)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
					return
				}

				offer := &Offer.Offer{Id: intid}
				offer, err = f.offerBase.GetOffer(offer)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find offer"})
					return
				}

				id := c.Param("consumerid")
				consumerid, err := strconv.Atoi(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
					return
				}

				if (curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != consumerid) || consumerid == offer.ProducerId {
					c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
					return
				}

				user := &User.User{Id:consumerid}
				user, err  = f.userBase.GetUser(user)
				if err != nil && errors.Is(err, pg.ErrNoRows) {
					c.JSON(http.StatusNotFound, gin.H{"No such id in system": consumerid})
					return
				}

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
					return
				}

				deal.FoodName = offer.FoodName
				deal.ConsumerId = user.Id
				deal.ConsumerName = user.UserName
				deal.ProducerId = offer.ProducerId
				deal.ProducerName = offer.ProducerName
				deal.Type = offer.Type
				deal.Myaso = offer.Myaso
				deal.Halal = offer.Halal
				deal.Vegan = offer.Vegan
				deal.Spicy = offer.Spicy
				deal.Description = offer.Description


				if deal.Quantity < 1 {
					c.JSON(http.StatusBadRequest, gin.H{"Error": "Wrong quantity"})
					return
				}
				if offer.AvailableQuantity < deal.Quantity {
					c.JSON(http.StatusBadRequest, gin.H{"Error": "too big quantity, not enough available"})
					return
				}

				offer.AvailableQuantity -= deal.Quantity
				if offer.AvailableQuantity != 0 {
					offer, err = f.offerBase.UpdateOffer(offer)

					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't update offer"})
						return
					}
				} else {
					offerlog := OfferLog.OfferLog(*offer)
					_, err = f.offerLogBase.CreateOfferLog(&offerlog)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't create offer log"})
						return
					}
					err = f.offerBase.DeleteOffer(intid)

					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't delete offer"})
						return
					}
				}

			}
		}
		deal.Created = time.Now()
		deal.Complete = "false"
		result, err := f.dealBase.CreateDeal(&deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't create deal"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func (f DealEndpointsFactory) ListDeals() func(c *gin.Context) {
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
		active := c.Query("onlyactive")
		var deals []*Deal.Deal
		idc := c.Param("consumerid")
		idp := c.Param("producerid")
		if (curruser.Permission == Authorization.Admin || curruser.Permission == Authorization.Manager) && len(idc) == 0 && len(idp) == 0 {
			if active == "true" {
				deals, err = f.dealBase.ListActiveDeals()
			} else {
				deals, err = f.dealBase.ListDeals()
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find deals"})
				return
			}
		} else {
			if len(idc) != 0 {
				intid, err := strconv.Atoi(idc)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
					return
				}

				if active == "true" {
					deals, err = f.dealBase.ListActiveConsumerDeals(intid)
				} else {
					deals, err = f.dealBase.ListConsumerDeals(intid)
				}

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find deals"})
					return
				}
			} else {
				if len(idp) != 0 {
					intid, err := strconv.Atoi(idp)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
						return
					}

					if active == "true" {
						deals, err = f.dealBase.ListActiveProducerDeals(intid)
					} else {
						deals, err = f.dealBase.ListProducerDeals(intid)
					}

					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find deals"})
						return
					}
				} else {
					c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
					return
				}
			}
		}
		c.JSON(http.StatusOK, deals)
	}
}

func (f DealEndpointsFactory) UpdateDeal() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}
		/*
			id := c.Param("consumerid")
			if len(id) == 0 {
				c.JSON(http.StatusBadRequest,gin.H{"Error": "No id provided"})
				return
			}

			userid, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error": "Provided id is not integer"})
				return
			}*/

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager /*&& curruser.UserId != userid*/ {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		dealtogetid := &Deal.Deal{Id: intid}
		dealtogetid, err = f.dealBase.GetDeal(dealtogetid)
		if err != nil && errors.Is(err, pg.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"No such id in system": intid})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Db Error"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager /*&& userid != dealtogetid.ProducerId && userid != dealtogetid.ConsumerId*/ {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		deal := &Deal.Deal{}
		err = c.ShouldBindJSON(&deal)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Provided data is in wrong format"})
			return
		}

		deal.Id = dealtogetid.Id
		deal.ProducerId = dealtogetid.ProducerId
		deal.ProducerName = dealtogetid.ProducerName
		deal.ConsumerId = dealtogetid.ConsumerId
		deal.ConsumerName = dealtogetid.ConsumerName
		deal.Created = dealtogetid.Created

		result, err := f.dealBase.UpdateDeal(deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't Update deal"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (f DealEndpointsFactory) CompleteDeal() func(c *gin.Context) {
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

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		deal := &Deal.Deal{Id: intid}
		deal, err = f.dealBase.GetDeal(deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find deal"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != deal.ProducerId {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		deal.Complete = "true"
		deal.Finished = time.Now()
		result, err := f.dealBase.UpdateDeal(deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't update deal"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (f DealEndpointsFactory) DeleteDeal() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager {
			c.JSON(http.StatusForbidden, gin.H{"Error": "Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Provided id is not integer"})
			return
		}

		err = f.dealBase.DeleteDeal(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't delete deal "})
			return
		}

		c.JSON(http.StatusOK, gin.H{"DealID": intid})
	}
}
