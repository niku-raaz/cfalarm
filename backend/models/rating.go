package models

import "time"

type Rating struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint
	Platform  string    `gorm:"size:20"`
	Rating    int
	Timestamp time.Time
}
