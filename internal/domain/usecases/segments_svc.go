// The usecases package implements the application's business logic. Since the functions are simple
// and there is almost no preliminary preparation before working with data, we immediately call the storage methods.
// It seems that testing can be neglected.
package usecases

import (
	"fmt"
	"segmentation-service/internal/domain/models"
	"segmentation-service/internal/ports"

	"github.com/google/uuid"
)

type SegmentSvc struct {
	storage ports.SegmentStorage
}

var _ ports.SegmentService = (*SegmentSvc)(nil)

// New returns a new instance of SegmentSvc.
func New(storage ports.SegmentStorage) *SegmentSvc {
	return &SegmentSvc{
		storage: storage,
	}
}

func (a *SegmentSvc) CreateSegment(slug string) error {
	count, err := a.storage.FindSegment(slug)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if count != 0 {
		return models.ErrSegmentAlreadyExists
	}

	err = a.storage.SaveSegment(slug)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (a *SegmentSvc) DeleteSegment(slug string) error {
	count, err := a.storage.FindSegment(slug)
	if err != nil {
		return err
	}
	if count == 0 {
		return models.ErrSegmentNotFound
	}
	err = a.storage.DeleteSegment(slug)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (a *SegmentSvc) UpdateUserSegments(data models.UpdateRequest, userID uuid.UUID) error {
	return a.storage.UpdateUserSegments(data, userID)
}

func (a *SegmentSvc) GetUserSegments(userID uuid.UUID) (models.SegmentsList, error) {
	return a.storage.GetUserSegments(userID)
}

func (a *SegmentSvc) GetReport(period string) ([][]string, error) {
	monthBeginning := period + "-01"
	records, err := a.storage.GetReport(monthBeginning)
	if err != nil {
		return records, err
	}
	return records, nil
}

func (a *SegmentSvc) GetUserReport(period string, userID uuid.UUID) ([][]string, error) {
	monthBeginning := period + "-01"
	records, err := a.storage.GetUserReport(monthBeginning, userID)
	if err != nil {
		return records, err
	}
	return records, nil
}
