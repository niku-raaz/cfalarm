package cron

import (
	"github.com/robfig/cron/v3"
	"log"
)

func Start() {
	c := cron.New()

	// Every day @ midnight
	_, _ = c.AddFunc("0 0 * * *", func() {
		RunDailyPractice()
	})

	// Every day @ 21:00
	_, _ = c.AddFunc("0 21 * * *", func() {
		RunDailyReminder()
	})

	// Every 2 days @ midnight
	_, _ = c.AddFunc("* * * * *", func() {
		RunAutoContestRegistration()
	})

	c.Start()
	log.Println("Cron scheduler started")
}
