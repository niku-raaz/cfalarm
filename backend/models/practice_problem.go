package models

import "time"

type PracticeProblem struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint
	ProblemID    string    `gorm:"size:50"`
	Rating       int
	DateAssigned time.Time `gorm:"type:date"`
	Status       string    `gorm:"size:20"`
	SolvedAt     *time.Time
}
