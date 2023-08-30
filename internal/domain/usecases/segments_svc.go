// The usecases package implements the application's business logic. Since the functions are simple
// and no preliminary preparation is required before working with data, we immediately call the storage methods.
package usecases

import (
	"context"
	"encoding/csv"
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

func (a *SegmentSvc) CreateSegment(ctx context.Context, slug string) error {
	count, err := a.storage.FindSegment(ctx, slug)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if count != 0 {
		return models.ErrSegmentAlreadyExists
	}

	err = a.storage.SaveSegment(ctx, slug)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (a *SegmentSvc) DeleteSegment(ctx context.Context, slug string) error {
	count, err := a.storage.FindSegment(ctx, slug)
	if err != nil {
		return err
	}
	if count == 0 {
		return models.ErrSegmentNotFound
	}
	err = a.storage.DeleteSegment(ctx, slug)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (a *SegmentSvc) UpdateUserSegments(ctx context.Context, data models.UpdateRequest, userID uuid.UUID) error {
	return a.storage.UpdateUserSegments(ctx, data, userID)
}

func (a *SegmentSvc) GetUserSegments(ctx context.Context, userID uuid.UUID) (models.SegmentsList, error) {
	return a.storage.GetUserSegments(ctx, userID)
}

func (a *SegmentSvc) GetReport(ctx context.Context, period string, wr *csv.Writer) error {
	monthBeginning := period + "-01"
	records, err := a.storage.GetReport(ctx, monthBeginning)
	if err != nil {
		return err
	}
	return wr.WriteAll(records)
}

func (a *SegmentSvc) GetUserReport(ctx context.Context, period string, userID uuid.UUID, wr *csv.Writer) error {
	monthBeginning := period + "-01"
	records, err := a.storage.GetUserReport(ctx, monthBeginning, userID)
	if err != nil {
		return err
	}
	return wr.WriteAll(records)
}
