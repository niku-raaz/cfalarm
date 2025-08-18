package models

import "time"

type Contest struct {
	ID              uint      `gorm:"primaryKey"`
	ContestID       int
	Name            string    `gorm:"size:150"`
	StartTime       time.Time
	DurationSeconds int
	CreatedAt       time.Time
}
