package models

import "time"

type EmailReminderLog struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint
	Date   time.Time `gorm:"type:date"`
	SentAt time.Time
}
