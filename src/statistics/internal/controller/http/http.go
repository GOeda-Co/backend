package http

import (
	"net/http"

	"repeatro/src/statistics/internal/service/statistics"
	"repeatro/src/statistics/pkg/scheme"

	"repeatro/internal/tools"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	ResultService *statistics.Service
}

func NewController(resultService *statistics.Service) *Controller {
	return &Controller{ResultService: resultService}
}

func (rc Controller) GetStats(ctx *gin.Context) {
	var interval schemes.Interval
	if err := ctx.ShouldBindJSON(&interval); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userId, err := tools.GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meanGrade, err := rc.ResultService.GetMeanGradeOfPeriod(interval.DtStart, interval.DtEnd, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"Mean Grade": meanGrade})
}
