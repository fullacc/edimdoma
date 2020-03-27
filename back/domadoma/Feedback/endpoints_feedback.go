package Feedback

import (
	"../Deal"
	"github.com/fullacc/edimdoma/back/domadoma"
	"github.com/fullacc/edimdoma/back/domadoma/User"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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
	userBase     User.UserBase
}

func (f FeedbackEndpointsFactory) GetFeedback() func(c *gin.Context) {
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

		id := c.Param( "feedbackid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		feedback, err := f.feedbackBase.GetFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,feedback)
	}
}

func (f FeedbackEndpointsFactory) CreateFeedback() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := domadoma.GetToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		id := c.Param( "consumerid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"No id provided"})
			return
		}

		consumerid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.UserId != consumerid {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		feedback := Feedback{}
		err = c.ShouldBindJSON(&feedback)
		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": err.Error()})
			return
		}

		id = c.Param( "dealid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ":"No id provided"})
			return
		}

		dealid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		}

		dealtogetid, err := f.dealBase.GetDeal(dealid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ":err.Error()})
			return
		}

		if curruser.Permission!= User.Admin && curruser.Permission!= User.Manager && consumerid != dealtogetid.ConsumerId{
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}

		feedback.ConsumerId = dealtogetid.ConsumerId
		feedback.ProducerId = dealtogetid.ProducerId
		feedback.Created = time.Now()
		feedback.DealId = dealtogetid.Id
		result, err := f.feedbackBase.CreateFeedback(&feedback)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error ": err.Error()})
			return
		}

		user := &User.User{Id:dealtogetid.ProducerId}
		user,_ = f.userBase.GetUser(user)
		user.RatingN ++
		user.RatingTotal += float64(feedback.Value)
		user.Rating = user.RatingTotal / user.RatingN
		_,_ = f.userBase.UpdateUser(user)


		c.JSON(http.StatusCreated,result)
	}
}

func (f FeedbackEndpointsFactory) ListFeedbacks() func(c *gin.Context) {
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

		var feedbacks []*Feedback
		idp := c.Param("producerid")
		if (curruser.Permission == User.Admin || curruser.Permission == User.Manager)&&len(idp) == 0 {
			feedbacks, err = f.feedbackBase.ListFeedbacks()
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

				feedbacks, err = f.feedbackBase.ListProducerFeedbacks(intid)
				if err != nil {
					c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
					return
				}

			} else {
					c.JSON(http.StatusForbidden,gin.H{"Error ": "Not allowed"})
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
			c.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.Permission != User.Regular {
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		id := c.Param("feedbackid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		feedbacktocheck, err := f.feedbackBase.GetFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		feedback := &Feedback{}
		err = c.ShouldBindJSON(feedback)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error ": err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.UserId != feedbacktocheck.ConsumerId{
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		feedback.ConsumerId = feedbacktocheck.ConsumerId
		feedback.ProducerId = feedbacktocheck.ProducerId
		feedback.DealId = feedbacktocheck.DealId
		if feedback.Created.IsZero() {
			feedback.Created = feedbacktocheck.Created
		}

		if feedback.Text == "" {
			feedback.Text = feedbacktocheck.Text
		}

		if feedback.Value == 0 {
			feedback.Value = feedbacktocheck.Value
		}

		user := &User.User{Id:curruser.UserId}
		user,_ = f.userBase.GetUser(user)
		user.RatingTotal += float64(feedback.Value) - float64(feedbacktocheck.Value)
		user.Rating = user.RatingTotal / user.RatingN
		_,_ = f.userBase.UpdateUser(user)

		feedback.Id = intid
		result, err := f.feedbackBase.UpdateFeedback(feedback)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func (f FeedbackEndpointsFactory) DeleteFeedback() func(c *gin.Context) {
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

		id := c.Param("feedbackid")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest,gin.H{"Error ": "No id provided"})
			return
		}

		intid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		feedbacktocheck, err := f.feedbackBase.GetFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		if curruser.Permission != User.Admin && curruser.Permission != User.Manager && curruser.UserId != feedbacktocheck.ConsumerId{
			c.JSON(http.StatusForbidden, gin.H{"Error ": "Not allowed"})
			return
		}

		user := &User.User{Id:curruser.UserId}
		user, _ = f.userBase.GetUser(user)
		user.RatingTotal -= float64(feedbacktocheck.Value)
		user.RatingN --
		_, _ = f.userBase.UpdateUser(user)

		err = f.feedbackBase.DeleteFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ": err.Error()})
			return
		}

		c.JSON(http.StatusOK,gin.H{"deletedid":intid})
	}
}
