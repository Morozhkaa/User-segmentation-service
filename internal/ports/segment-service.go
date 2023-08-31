// The ports package contains the description of the interfaces.
package ports

import (
	"segmentation-service/internal/domain/models"

	"github.com/google/uuid"
)

type SegmentService interface {
	CreateSegment(slug string) error
	DeleteSegment(slug string) error
	UpdateUserSegments(data models.UpdateRequest, userID uuid.UUID) error
	GetUserSegments(userID uuid.UUID) (models.SegmentsList, error)
	GetReport(period string) ([][]string, error)
	GetUserReport(period string, userID uuid.UUID) ([][]string, error)
}
