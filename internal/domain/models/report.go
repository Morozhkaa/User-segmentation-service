package models

import (
	"time"

	"github.com/google/uuid"
)

type reportAction string

const (
	ActAdd    = "add"
	ActRemove = "remove"
)

type ReportRow struct {
	UserID      uuid.UUID
	SegmentName string
	Action      reportAction
	Time        time.Time
}
