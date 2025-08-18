package models

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey"`
	Name          string    `gorm:"size:100"`
	Email         string    `gorm:"size:150;unique"`
	PasswordHash  string    `gorm:"type:text"`
	CodeforcesID  string    `gorm:"size:50"`
	CodechefID    string    `gorm:"size:50"`
	AtcoderID     string    `gorm:"size:50"`
	LeetcodeID    string    `gorm:"size:50"`
	GoogleToken   string    `gorm:"type:text"`   // <-- NEW
	CreatedAt     time.Time
}
