package cron

import (
	"log"
	"time"

	"cfalarm/config"
	"cfalarm/models"
	"cfalarm/services"
)

func RunDailyReminder() {
	log.Println("[cron] Checking for unsolved problems...")

	var users []models.User
	config.DB.Find(&users)

	today := time.Now().Format("2006-01-02")

	for _, user := range users {
		var count int64
		config.DB.Model(&models.PracticeProblem{}).
			Where("user_id = ? AND date_assigned = ? AND status = 'unsolved'", user.ID, today).
			Count(&count)

		if count > 0 {
			// check if a reminder was already sent
			var logEntry models.EmailReminderLog
			if err := config.DB.Where("user_id = ? AND date = ?", user.ID, today).
				First(&logEntry).Error; err == nil {
				continue
			}

			// send email
			err := services.SendReminderEmail(user.Email, "Reminder: You have unsolved practice problems for today!")
			if err == nil {
				// log it
				config.DB.Create(&models.EmailReminderLog{
					UserID: user.ID,
					Date:   time.Now(),
				})
			}
		}
	}
}
