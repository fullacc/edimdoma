package domadoma

import (
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
}

func (f FeedbackEndpointsFactory) GetFeedback() func(c *gin.Context) {
	return func(c *gin.Context) {
		curruser, err := getToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}
		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.Permission != Regular{
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
		curruser, err := getToken(c.Request.Header.Get("Token"))
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error :":err.Error()})
			return
		}
		if curruser.Permission != Admin && curruser.Permission != Manager && curruser.Permission != Regular{
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
		dealtogetid, err:= postgreBase.GetDeal(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error ":err.Error()})
			return
		}
		if curruser.Permission!=Admin && curruser.Permission!=Manager && curruser.UserId != dealtogetid.ConsumerId{
			c.JSON(http.StatusForbidden,gin.H{"Error":"Not allowed"})
			return
		}
		var feedback *Feedback
		if err := c.ShouldBindJSON(&feedback); err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest,gin.H{"Error: ": err.Error()})
				return
			}
		}
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
		CHECKIFAUTHORIZED
		var feedbacks []*Feedback
		feedbacks, err := f.feedbackBase.ListFeedbacks()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusCreated,feedbacks)
	}
}

func (f FeedbackEndpointsFactory) UpdateFeedback() func(c *gin.Context) {
	return func(c *gin.Context) {
		CHECKAUTHORIZED
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
		feedback := &Feedback{}
		if err := c.ShouldBindJSON(feedback); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
			return
		}
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
		CHECKAUTHORIZED
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
		err = f.feedbackBase.DeleteFeedback(intid)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"Error: ": err.Error()})
			return
		}
		c.JSON(http.StatusOK,gin.H{"deleted":intid})
	}
}
