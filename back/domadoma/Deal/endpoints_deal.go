package Deal

import (
	"github.com/fullacc/edimdoma/back/domadoma"
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

func NewDealEndpoints(dealBase DealBase) DealEndpoints {
	return &DealEndpointsFactory{dealBase: dealBase}
}

type DealEndpointsFactory struct{
	dealBase DealBase
	offerBase Offer.OfferBase
	offerLogBase OfferLog.OfferLogBase
	requestBase Request.RequestBase
}

func (d DealEndpointsFactory) GetDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden,gin.H{"Error: ":"Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		deal, err := d.dealBase.GetDeal(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.UserId != deal.ConsumerId && curruser.UserId != deal.ProducerId {
			c.JSON(http.StatusForbidden,gin.H{"Error :": "Not allowed"})
			return
		}

		c.JSON(http.StatusOK,deal)
	}
}

func (d DealEndpointsFactory) CreateDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		deal := Deal{}
		if err := c.ShouldBindJSON(&deal); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		reqid := c.Param("requestid")
		offid := c.Param("offerid")
		if len(reqid) != 0 {
			intid, err := strconv.Atoi(reqid)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
				return
			}

			err = d.requestBase.DeleteRequest(intid)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error: ":err.Error()})
				return
			}

		} else {
			if len(offid) != 0 {
				intid, err := strconv.Atoi(offid)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
					return
				}

				offer, err := d.offerBase.GetOffer(intid)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
					return
				}

				offer.AvailableQuantity -= deal.Quantity
				if offer.AvailableQuantity >0 {
					offer, err = d.offerBase.UpdateOffer(offer)
				} else {
					offerlog := OfferLog.OfferLog(*offer)
					_,_ = d.offerLogBase.CreateOfferLog(&offerlog)
					err = d.offerBase.DeleteOffer(intid)
				}

				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error: ":err.Error()})
					return
				}
			}
		}

		result, err := d.dealBase.CreateDeal(&deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		c.JSON(http.StatusCreated,result)
	}
}

func (d DealEndpointsFactory) ListDeals() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}
		active := c.Query("onlyactive")
		var deals []*Deal
		idc := c.Param( "consumerid")
		idp := c.Param("producerid")
		if (curruser.Permission == domadoma.Admin || curruser.Permission == domadoma.Manager)&& len(idc)==0 && len(idp) == 0 {
			if active == "true"{
				deals, err = d.dealBase.ListActiveDeals()
			} else {
				deals, err = d.dealBase.ListDeals()
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
				return
			}
		} else{
			if len(idc) != 0 {
				intid, err := strconv.Atoi(idc)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
					return
				}

				if active == "true" {
					deals, err = d.dealBase.ListActiveConsumerDeals(intid)
				} else {
					deals, err = d.dealBase.ListConsumerDeals(intid)
				}

				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
					return
				}
			} else {
				if len(idp) != 0 {
					intid, err := strconv.Atoi(idp)
					if err != nil {
						c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
						return
					}

					if active == "true" {
						deals,err = d.dealBase.ListActiveProducerDeals(intid)
					} else {
						deals, err = d.dealBase.ListProducerDeals(intid)
					}

					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
						return
					}
				} else {
					c.JSON(http.StatusForbidden,gin.H{"Error: ": "Not allowed"})
					return
				}
			}
		}
		c.JSON(http.StatusOK,deals)
	}
}


func (d DealEndpointsFactory) UpdateDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err!=nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		dealtogetid, err := d.dealBase.GetDeal(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ":err.Error()})
			return
		}

		if curruser.Permission!= domadoma.Admin && curruser.Permission!= domadoma.Manager && curruser.UserId != dealtogetid.ProducerId && curruser.UserId != dealtogetid.ConsumerId{
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}

		deal := Deal{}
		if err := c.ShouldBindJSON(&deal); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		deal.Id = intid
		result, err := d.dealBase.UpdateDeal(&deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		c.JSON(http.StatusOK,result)
	}
}

func (d DealEndpointsFactory) CompleteDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.Permission != domadoma.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err!=nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		deal, err := d.dealBase.GetDeal(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ":err.Error()})
			return
		}

		if curruser.Permission!= domadoma.Admin && curruser.Permission!= domadoma.Manager && curruser.UserId != deal.ProducerId {
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}

		deal.Complete = true
		deal.Finished = time.Now()
		result, err := d.dealBase.UpdateDeal(deal)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		c.JSON(http.StatusOK,result)
	}
}

func (d DealEndpointsFactory) DeleteDeal() func(c *gin.Context) {
	return func(c *gin.Context){
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager {
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		id := c.Param("dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id given"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = d.dealBase.DeleteDeal(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"DealID": intid})
	}
}
