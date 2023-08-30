package http

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"regexp"
	"segmentation-service/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @ID createSegment
// @tags segment
// @Summary Create a new segment
// @Description Creates a new segment with the given slug. If this segment was already in the database, return the BadRequest status.
// @Accept json
// @Param slug body models.Segment true "A short name containing only letters, numbers, underscores, or hyphens. Format: ^[\w-]+$"
// @Success 201 {object} models.SuccessResponse "Segment created successfully."
// @Failure 400 {object} models.ErrorResponse "Segment already exists / missing required 'slug' parameter / invalid format of 'slug' parameter."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /createSegment [post]
func (a *Adapter) createSegment(ctx *gin.Context) {
	var segment models.Segment
	err := ctx.BindJSON(&segment)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	matched, err := regexp.MatchString(`^[\w-]+$`, string(segment.Slug))
	if err != nil || !matched {
		a.ErrorHandler(ctx, models.ErrInvalidSlugFormat)
		return
	}

	err = a.segmentSvc.CreateSegment(ctx, segment.Slug)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}

	ctx.IndentedJSON(
		http.StatusCreated,
		models.SuccessResponse{SuccessMsg: fmt.Sprintf("segment with slug '%s' created", segment.Slug)},
	)
}

// @ID deleteSegment
// @tags segment
// @Summary Delete segment
// @Description Delete the segment with the given slug and all users from it.
// @Accept json
// @Param slug body models.Segment true "A short name containing only letters, numbers, underscores, or hyphens. Format: ^[\w-]+$"
// @Success 200 {object} models.SuccessResponse "Segment deleted successfully."
// @Failure 400 {object} models.ErrorResponse "Missing required 'slug' parameter / invalid format of 'slug' parameter."
// @Failure 404 {object} models.ErrorResponse "Segment with the given slug not found."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /deleteSegment [delete]
func (a *Adapter) deleteSegment(ctx *gin.Context) {
	var segment models.Segment
	err := ctx.BindJSON(&segment)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	matched, err := regexp.MatchString(`^[\w-]+$`, string(segment.Slug))
	if err != nil || !matched {
		a.ErrorHandler(ctx, models.ErrInvalidSlugFormat)
		return
	}

	err = a.segmentSvc.DeleteSegment(ctx, segment.Slug)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}

	ctx.IndentedJSON(
		http.StatusOK,
		models.SuccessResponse{SuccessMsg: fmt.Sprintf("segment with slug '%s' deleted", segment.Slug)},
	)
}

// @ID updateSegments
// @tags segment
// @Summary Update user segments
// @Description Add/remove a user from segments in accordance with the transferred lists for adding and deleting.
// @Accept json
// @Param userID path string true "User ID in uuid format" Format(uuid)
// @Param segments body models.UpdateRequest true "segments"
// @Success 200 {object} models.SuccessResponse "User information updated successfully."
// @Failure 400 {object} models.ErrorResponse "Invalid format for parameter 'userID'."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /updateUserSegments/{userID} [post]
func (a *Adapter) updateSegments(ctx *gin.Context) {
	user_id, err := a.getIdFromPath(ctx)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	var data models.UpdateRequest
	err = ctx.BindJSON(&data)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}

	err = a.segmentSvc.UpdateUserSegments(ctx, data, user_id)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.IndentedJSON(
		http.StatusOK,
		models.SuccessResponse{SuccessMsg: fmt.Sprintf("segment information for user with userID = %v updated", user_id)},
	)
}

// @ID getSegments
// @tags segment
// @Summary Get user segments
// @Description Return the list of segments the user is a member of.
// @Param userID path string true "User ID in uuid format" Format(uuid)
// @Success 200 {object} models.SegmentsList "User segments received successfully."
// @Failure 400 {object} models.ErrorResponse "Invalid format for parameter 'userID'."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /getUserSegments/{userID} [get]
func (a *Adapter) getSegments(ctx *gin.Context) {
	user_id, err := a.getIdFromPath(ctx)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	segments, err := a.segmentSvc.GetUserSegments(ctx, user_id)
	if err != nil {
		a.ErrorHandler(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	ctx.IndentedJSON(http.StatusOK, segments)
}

// @ID getReport
// @tags report
// @Summary Get report file
// @Description Returns the history of events for the given month as a csv file.
// @Accept json
// @Produce text/csv
// @Param period path string true "Month for which you want to display information, in the format 'yyyy-mm'"
// @Success 200 "Report file received successfully."
// @Failure 400 {object} models.ErrorResponse "Invalid format for parameter 'period'."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /getReport/{period} [get]
func (a *Adapter) getReport(ctx *gin.Context) {
	period, err := a.getPeriodFromPath(ctx)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}

	// set multiple http headers so that the browser responds by downloading the CSV file
	ctx.Writer.Header().Set("Content-Type", "text/csv")
	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename=data.csv")
	wr := csv.NewWriter(ctx.Writer)

	err = a.segmentSvc.GetReport(ctx, period, wr)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
}

// @ID getUserReport
// @tags report
// @Summary Get a report file for a specific user
// @Description Returns a specific user's history of events for the specified month as a csv file.
// @Accept json
// @Produce text/csv
// @Param period path string true "Month for which you want to display information, in the format 'yyyy-mm'"
// @Param userID path string true "User ID in uuid format" Format(uuid)
// @Success 200 "Report file received successfully."
// @Failure 400 {object} models.ErrorResponse "Invalid format for parameter 'period'."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /getUserReport/{period}/{userID} [get]
func (a *Adapter) getUserReport(ctx *gin.Context) {
	period, err := a.getPeriodFromPath(ctx)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	user_id, err := a.getIdFromPath(ctx)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}

	// set multiple http headers so that the browser responds by downloading the CSV file
	ctx.Writer.Header().Set("Content-Type", "text/csv")
	ctx.Writer.Header().Set("Content-Disposition", "attachment;filename=userdata.csv")
	wr := csv.NewWriter(ctx.Writer)

	err = a.segmentSvc.GetUserReport(ctx, period, user_id, wr)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
}

func (a *Adapter) getIdFromPath(ctx *gin.Context) (uuid.UUID, error) {
	if ctx.Param("userID") == ":userID" { // if path-parameter is not set (this is how it works in Postman)
		return uuid.Nil, models.ErrBadRequest
	}
	user_id, err := uuid.Parse(ctx.Param("userID")) // may be invalid UUID length or formatis
	if err != nil {
		return uuid.Nil, models.ErrInvalidUuidFormat
	}
	return user_id, nil
}

func (a *Adapter) getPeriodFromPath(ctx *gin.Context) (string, error) {
	period := ctx.Param("period")
	if period == ":period" {
		return "", models.ErrBadRequest
	}
	matched, err := regexp.MatchString(`^\d{4}-\d{2}$`, period)
	if err != nil || !matched {
		return "", models.ErrInvalidPeriodFormat
	}
	return period, nil
}
