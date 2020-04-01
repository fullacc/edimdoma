package Endpoints

import (
	"github.com/fullacc/edimdoma/back/domadoma/Authorization"
	"github.com/fullacc/edimdoma/back/domadoma/Deal"
	"github.com/fullacc/edimdoma/back/domadoma/Offer"
	"github.com/fullacc/edimdoma/back/domadoma/OfferLog"
	"github.com/fullacc/edimdoma/back/domadoma/Request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type DealEndpoints interface{
	CreateDeal() func(c *gin.Context)

	GetDeal() func(c *gin.Context)

	ListDeals() func(c *gin.Context)

	UpdateDeal() func(c *gin.Context)

	DeleteDeal() func(c *gin.Context)

	CompleteDeal() func(c *gin.Context)

}

func NewDealEndpoints(dealBase Deal.DealBase, authorizationBase Authorization.AuthorizationBase, offerBase Offer.OfferBase, offerLogBase OfferLog.OfferLogBase, requestBase Request.RequestBase) DealEndpoints {
	return &DealEndpointsFactory{dealBase: dealBase, authorizationBase:authorizationBase, offerBase:offerBase, offerLogBase:offerLogBase, requestBase:requestBase}
}

type DealEndpointsFactory struct{
	authorizationBase Authorization.AuthorizationBase
	dealBase          Deal.DealBase
	offerBase         Offer.OfferBase
	offerLogBase      OfferLog.OfferLogBase
	requestBase       Request.RequestBase
}

func (f DealEndpointsFactory) GetDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden,gin.H{"Error ":"Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}

		dl := Deal.Deal{Id:intid}
		deal, err := f.dealBase.GetDeal(&dl)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find deal"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != deal.ConsumerId && curruser.UserId != deal.ProducerId {
			c.JSON(http.StatusForbidden,gin.H{"Error": "Not allowed"})
			return
		}

		c.JSON(http.StatusOK,deal)
	}
}

func (f DealEndpointsFactory) CreateDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		reqid := c.Param("requestid")
		offid := c.Param("offerid")

		deal := Deal.Deal{}
		err = c.ShouldBindJSON(&deal)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provided data is in wrong format"})
			return
		}

		if deal.Quantity < 1 {
			c.JSON(http.StatusBadRequest,gin.H{"Error":"Wrong quantity"})
			return
		}

		if deal.Price < 1{
			c.JSON(http.StatusBadRequest,gin.H{"Error":"Wrong price"})
			return
		}

		if len(reqid) != 0 {
			intid, err := strconv.Atoi(reqid)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
				return
			}

			rq := Request.Request{Id:intid}
			request, err := f.requestBase.GetRequest(&rq)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find request"})
				return
			}

			id := c.Param("producerid")
			producerid, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
				return
			}

			if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != producerid {
				c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
				return
			}

			deal.Quantity = request.Quantity
			deal.Food = request.Food
			deal.ConsumerId = request.ConsumerId
			deal.ProducerId = producerid

			err = f.requestBase.DeleteRequest(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error ":"Couldn't delete request"})
				return
			}

		} else {
			if len(offid) != 0 {
				intid, err := strconv.Atoi(offid)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
					return
				}

				ofr := Offer.Offer{Id:intid}
				offer, err := f.offerBase.GetOffer(&ofr)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find offer"})
					return
				}

				id := c.Param("consumerid")
				consumerid, err := strconv.Atoi(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
					return
				}

				if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.UserId != consumerid {
					c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
					return
				}

				deal.Food = offer.Food
				deal.ConsumerId = consumerid
				deal.ProducerId = offer.ProducerId

				if offer.AvailableQuantity < deal.Quantity {
					c.JSON(http.StatusBadRequest,gin.H{"Error":"too big quantity, not enough available"})
					return
				}

				offer.AvailableQuantity -= deal.Quantity
				if offer.AvailableQuantity != 0 {
					offer, err = f.offerBase.UpdateOffer(offer)

					if err != nil {
						c.JSON(http.StatusInternalServerError,gin.H{"Error ":"Couldn't update offer"})
						return
					}
				} else {
					offerlog := OfferLog.OfferLog(*offer)
					_,err = f.offerLogBase.CreateOfferLog(&offerlog)
					if err != nil {
						c.JSON(http.StatusInternalServerError,gin.H{"Error ":"Couldn't create offer log"})
						return
					}
					err = f.offerBase.DeleteOffer(intid)

					if err != nil {
						c.JSON(http.StatusInternalServerError,gin.H{"Error ":"Couldn't delete offer"})
						return
					}
				}

			}
		}
		deal.Created = time.Now()
		deal.Complete = false
		result, err := f.dealBase.CreateDeal(&deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't create deal"})
			return
		}

		c.JSON(http.StatusCreated,result)
	}
}

func (f DealEndpointsFactory) ListDeals() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}
		active := c.Query("onlyactive")
		var deals []*Deal.Deal
		idc := c.Param( "consumerid")
		idp := c.Param("producerid")
		if (curruser.Permission == Authorization.Admin || curruser.Permission == Authorization.Manager)&& len(idc)==0 && len(idp) == 0 {
			if active == "true"{
				deals, err = f.dealBase.ListActiveDeals()
			} else {
				deals, err = f.dealBase.ListDeals()
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find deals"})
				return
			}
		} else{
			if len(idc) != 0 {
				intid, err := strconv.Atoi(idc)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
					return
				}

				if active == "true" {
					deals, err = f.dealBase.ListActiveConsumerDeals(intid)
				} else {
					deals, err = f.dealBase.ListConsumerDeals(intid)
				}

				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't find deals"})
					return
				}
			} else {
				if len(idp) != 0 {
					intid, err := strconv.Atoi(idp)
					if err != nil {
						c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
						return
					}

					if active == "true" {
						deals,err = f.dealBase.ListActiveProducerDeals(intid)
					} else {
						deals, err = f.dealBase.ListProducerDeals(intid)
					}

					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"Error ": "Couldn't find deals"})
						return
					}
				} else {
					c.JSON(http.StatusForbidden,gin.H{"Error ": "Not allowed"})
					return
				}
			}
		}
		c.JSON(http.StatusOK,deals)
	}
}


func (f DealEndpointsFactory) UpdateDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}
/*
		id := c.Param("consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		userid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Provided id is not integer"})
			return
		}*/

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager /*&& curruser.UserId != userid*/{
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err!=nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Provided id is not integer"})
			return
		}

		dl := Deal.Deal{Id:intid}
		dealtogetid, err := f.dealBase.GetDeal(&dl)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ":"Couldn't find deal"})
			return
		}

		if curruser.Permission!= Authorization.Admin && curruser.Permission!= Authorization.Manager /*&& userid != dealtogetid.ProducerId && userid != dealtogetid.ConsumerId*/{
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}

		deal := &Deal.Deal{}
		err = c.ShouldBindJSON(&deal)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provided data is in wrong format"})
			return
		}

		deal.Id = dealtogetid.Id
		deal.ProducerId = dealtogetid.ProducerId
		deal.ConsumerId = dealtogetid.ConsumerId
		deal.Created = dealtogetid.Created

		result, err := f.dealBase.UpdateDeal(deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't update deal"})
			return
		}

		c.JSON(http.StatusOK,result)
	}
}

func (f DealEndpointsFactory) CompleteDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager && curruser.Permission != Authorization.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err!=nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Provided id is not integer"})
			return
		}

		dl := Deal.Deal{Id:intid}
		deal, err := f.dealBase.GetDeal(&dl)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ":"Couldn't find deal"})
			return
		}

		if curruser.Permission!= Authorization.Admin && curruser.Permission!= Authorization.Manager && curruser.UserId != deal.ProducerId {
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}

		deal.Complete = true
		deal.Finished = time.Now()
		result, err := f.dealBase.UpdateDeal(deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't update deal"})
			return
		}

		c.JSON(http.StatusOK,result)
	}
}

func (f DealEndpointsFactory) DeleteDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := f.authorizationBase.GetAuthToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Couldn't find token"})
			return
		}

		if curruser.Permission != Authorization.Admin && curruser.Permission != Authorization.Manager {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Provided id is not integer"})
			return
		}

		err = f.dealBase.DeleteDeal(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": "Couldn't delete deal "})
			return
		}

		c.JSON(http.StatusOK,gin.H{"DealID": intid})
	}
}