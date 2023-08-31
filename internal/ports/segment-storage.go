package ports

import (
	"segmentation-service/internal/domain/models"

	"github.com/google/uuid"
)

type SegmentStorage interface {
	FindSegment(slug string) (int, error)
	SaveSegment(slug string) error
	DeleteSegment(slug string) error
	UpdateUserSegments(data models.UpdateRequest, userID uuid.UUID) error
	GetUserSegments(userID uuid.UUID) (models.SegmentsList, error)
	GetReport(period string) ([][]string, error)
	GetUserReport(period string, userID uuid.UUID) ([][]string, error)
}
