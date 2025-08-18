package models

import "time"

type UserContestRegistration struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint
	ContestID      uint
	RegisteredAt   time.Time
	CalendarEventID string  `gorm:"size:200"`
}
