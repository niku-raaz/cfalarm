package models

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey"`
	Name          string    `gorm:"size:100"`
	Email         string    `gorm:"size:150;unique"`
	CodeforcesID  string    `gorm:"size:50"`
	GoogleToken   string    `gorm:"type:text"`   
	CreatedAt     time.Time
    
}
 