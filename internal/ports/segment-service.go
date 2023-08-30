// The ports package contains the description of the interfaces.
package ports

import (
	"context"
	"encoding/csv"
	"segmentation-service/internal/domain/models"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SegmentService
type SegmentService interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	UpdateUserSegments(ctx context.Context, data models.UpdateRequest, userID uuid.UUID) error
	GetUserSegments(ctx context.Context, userID uuid.UUID) (models.SegmentsList, error)
	GetReport(ctx context.Context, period string, wr *csv.Writer) error
	GetUserReport(ctx context.Context, period string, userID uuid.UUID, wr *csv.Writer) error
}
