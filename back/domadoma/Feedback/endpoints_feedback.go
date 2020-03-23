package Feedback

import (
	"../Deal"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FeedbackEndpoints interface{
	GetFeedback() func(c *gin.Context)

	CreateFeedback() func(c *gin.Context)

	ListFeedbacks() func(c *gin.Context)

	UpdateFeedback() func(c *gin.Context)

	DeleteFeedback() func(c *gin.Context)

}

func NewFeedbackEndpoints(feedbackBase FeedbackBase) FeedbackEndpoints {
	return &FeedbackEndpointsFactory{feedbackBase: feedbackBase}
}

type FeedbackEndpointsFactory struct{
	feedbackBase FeedbackBase
	dealBase     Deal.DealBase
}

func (f FeedbackEndpointsFactory) GetFeedback() func(c *gin.Context) {
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

		id := c.Param( "feedbackid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		feedback, err := f.feedbackBase.GetFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,feedback)
	}
}

func (f FeedbackEndpointsFactory) CreateFeedback() func(c *gin.Context) {
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

		dealtogetid, err := f.dealBase.GetDeal(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ":err.Error()})
			return
		}

		if curruser.Permission!= domadoma.Admin && curruser.Permission!= domadoma.Manager && curruser.UserId != dealtogetid.ConsumerId{
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}

		var feedback *Feedback
		if err := c.ShouldBindJSON(&feedback); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": err.Error()})
			return
		}

		feedback.ConsumerId = dealtogetid.ConsumerId
		feedback.ProducerId = dealtogetid.ProducerId
		feedback.DealId = intid
		result, err := f.feedbackBase.CreateFeedback(feedback)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusCreated,result)
	}
}

func (f FeedbackEndpointsFactory) ListFeedbacks() func(c *gin.Context) {
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

		var feedbacks []*Feedback
		idp := c.Param("producerid")
		if (curruser.Permission == domadoma.Admin || curruser.Permission == domadoma.Manager)&&len(idp) == 0 {
			feedbacks, err = f.feedbackBase.ListFeedbacks()
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
				return
			}

		} else{
			if len(idp) != 0 {
				intid, err := strconv.Atoi(idp)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
					return
				}

				feedbacks, err = f.feedbackBase.ListProducerFeedbacks(intid)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
					return
				}

			} else {
					c.JSON(http.StatusForbidden,gin.H{"Error: ": "Not allowed"})
					return
				}
			}
		c.JSON(http.StatusOK,feedbacks)
	}
}

func (f FeedbackEndpointsFactory) UpdateFeedback() func(c *gin.Context) {
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

		id := c.Param("feedbackid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		feedbacktocheck, err := f.feedbackBase.GetFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		feedback := &Feedback{}
		if err := c.ShouldBindJSON(feedback); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.UserId != feedbacktocheck.ConsumerId{
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		feedback.ConsumerId = feedbacktocheck.ConsumerId
		feedback.ProducerId = feedbacktocheck.ProducerId
		feedback.DealId = feedbacktocheck.DealId
		feedback, err = f.feedbackBase.UpdateFeedback(intid, feedback)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,feedback)
	}
}

func (f FeedbackEndpointsFactory) DeleteFeedback() func(c *gin.Context) {
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

		id := c.Param("feedbackid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error: ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		feedbacktocheck, err := f.feedbackBase.GetFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		if curruser.Permission != domadoma.Admin && curruser.Permission != domadoma.Manager && curruser.UserId != feedbacktocheck.ConsumerId{
			c.JSON(http.StatusForbidden, gin.H{"Error: ": "Not allowed"})
			return
		}

		err = f.feedbackBase.DeleteFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deletedid":intid})
	}
}
