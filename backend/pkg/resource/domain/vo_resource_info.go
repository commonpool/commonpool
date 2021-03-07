package domain

import "time"

type ResourceValueType string

const (
	FromToDuration ResourceValueType = "from_to_duration"
)

type ResourceValueEstimation struct {
	ValueType         ResourceValueType `json:"value_type" gorm:"not null;type:varchar(128)"`
	ValueFromDuration time.Duration     `json:"time_value_from" gorm:"not null"`
	ValueToDuration   time.Duration     `json:"time_value_to" gorm:"not null"`
}

type ResourceInfo struct {
	Value       ResourceValueEstimation `json:"value"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
}
