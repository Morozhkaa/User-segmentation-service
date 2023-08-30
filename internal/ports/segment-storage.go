package ports

import (
	"context"
	"segmentation-service/internal/domain/models"

	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SegmentStorage
type SegmentStorage interface {
	FindSegment(ctx context.Context, slug string) (int, error)
	SaveSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	UpdateUserSegments(ctx context.Context, data models.UpdateRequest, userID uuid.UUID) error
	GetUserSegments(ctx context.Context, userID uuid.UUID) (models.SegmentsList, error)
	GetReport(ctx context.Context, period string) ([][]string, error)
	GetUserReport(ctx context.Context, period string, userID uuid.UUID) ([][]string, error)
}
