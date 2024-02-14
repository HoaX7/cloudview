package models

import (
	"time"

	"github.com/google/uuid"
)

// Enum (CPU_USAGE, RAM_USAGE, DISK_OPERATIONS)
type MetricType string
type TimeSeries struct {
	ID                 int32        `sql:"primary_key" json:"id"`
	TruncatedTimestamp int64        `sql:"primary_key" json:"truncatedTimestamp"`
	Series             *[][]float64 `sql:"series,double precision[][]" json:"series,omitempty"`
	Type               MetricType   `json:"type"`
	MetricPanelID      uuid.UUID    `json:"metricPanelId"`
	CreatedAt          *time.Time   `json:"createdAt,omitempty"`
	UpdatedAt          *time.Time   `json:"updatedAt,omitempty"`
}
