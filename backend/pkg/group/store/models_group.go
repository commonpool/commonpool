package store

import "time"

type Group struct {
	ID          string    `mapstructure:"id"`
	CreatedAt   time.Time `mapstructure:"createdAt"`
	Description string    `mapstructure:"description"`
	CreatedBy   string    `mapstructure:"createdBy"`
	Name        string    `mapstructure:"name"`
}
