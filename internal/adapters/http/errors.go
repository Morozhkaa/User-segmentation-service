package http

import (
	"errors"
	"net/http"
	"segmentation-service/internal/domain/models"
	"segmentation-service/pkg/infra/logger"

	"github.com/gin-gonic/gin"
)

func (a *Adapter) ErrorHandler(ctx *gin.Context, err error) {
	logger.Get().Warn("request failed: ", "desc", err.Error())

	switch {
	case errors.Is(err, models.ErrInvalidSlugFormat), errors.Is(err, models.ErrInvalidUuidFormat),
		errors.Is(err, models.ErrInvalidPeriodFormat), errors.Is(err, models.ErrBadRequest),
		errors.Is(err, models.ErrSegmentAlreadyExists):
		ctx.JSON(
			http.StatusBadRequest,
			models.ErrorResponse{ErrorMsg: err.Error()},
		)
	case errors.Is(err, models.ErrSegmentNotFound):
		ctx.JSON(
			http.StatusNotFound,
			models.ErrorResponse{ErrorMsg: err.Error()},
		)
	default:
		ctx.JSON(
			http.StatusInternalServerError,
			models.ErrorResponse{ErrorMsg: err.Error()},
		)
	}
}
