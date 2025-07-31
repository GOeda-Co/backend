package http

import (
	"fmt"
	"net/http"

	statsv1 "github.com/GOeda-Co/proto-contract/gen/go/stats"
	"github.com/gin-gonic/gin"
	_ "github.com/tomatoCoderq/repeatro/pkg/models"
)

// GetCardsReviewedCount godoc
// @Summary      Get user's cards reviewed count
// @Description  Returns the count of cards reviewed by the current user for the daily time range
// @Tags         statistics
// @Produce      json
// @Success      200  {object}  statsv1.GetCardsReviewedCountResponse
// @Failure      400  {object}  model.ErrorResponse	"Bad Request - Failed to get user ID from context"
// @Failure      500  {object}  model.ErrorResponse	"Internal Server Error - Failed to get reviewed cards"
// @Router       /stats/count [get]
func (cc *Controller) GetCardsReviewedCount(ctx *gin.Context) {
	// did := ctx.Param("id")

	uid, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error during getting uid occurred: %v", err)})
		return
	}

	response, err := cc.statClient.GetCardsReviewedCount(ctx, uid.String(), "", statsv1.TimeRange_DAILY)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get reviewed cards: %v", err)})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// GetAverageGrade godoc
// @Summary      Get user's average grade
// @Description  Returns the average grade for the current user for the daily time range
// @Tags         statistics
// @Produce      json
// @Success      200  {object}  statsv1.GetAverageGradeResponse
// @Failure      400  {object}  model.ErrorResponse	"Bad Request - Failed to get user ID from context"
// @Failure      500  {object}  model.ErrorResponse	"Internal Server Error - Failed to get average grade"
// @Router       /stats/average [get]
func (cc *Controller) GetAverageGrade(ctx *gin.Context) {
	// did := ctx.Param("id")

	uid, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := cc.statClient.GetAverageGrade(ctx, uid.String(), "", statsv1.TimeRange_DAILY)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get average grade: %v", err)})
		return
	}
	ctx.JSON(http.StatusOK, response)
}