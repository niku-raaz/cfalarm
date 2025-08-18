package cron

import (
	"context"
	"log"
	"time"

	"cfalarm/config"
	"cfalarm/models"
	"cfalarm/services"
	"encoding/json"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func RunAutoContestRegistration() {
	log.Println("[cron] Auto contest registration...")

	// 1. Get upcoming CF contests (placeholder)
	// Example contest (replace later with CF API response)
	startTime := time.Now().Add(48 * time.Hour) // two days from now
	endTime := startTime.Add(2 * time.Hour)     // 2h duration
	contestName := "Codeforces Round (Example)"

	var users []models.User
	config.DB.Find(&users)

	for _, user := range users {
		// Insert a record in user_contest_registrations
		reg := models.UserContestRegistration{
			UserID:       user.ID,
			ContestID:    0,
			RegisteredAt: time.Now(),
		}
		config.DB.Create(&reg)

		// Create Google Calendar event
		var token oauth2.Token
		_ = json.Unmarshal([]byte(user.GoogleToken), &token)

		ctx := context.Background()
		srv, err := calendar.NewService(ctx, option.WithTokenSource(services.GetGoogleOAuthConfig().TokenSource(ctx, &token)))
		if err != nil {
			log.Println("Calendar service error:", err)
			continue
		}

		event := &calendar.Event{
			Summary: contestName,
			Start: &calendar.EventDateTime{
				DateTime: startTime.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: endTime.Format(time.RFC3339),
			},
		}
		created, err := srv.Events.Insert("primary", event).Do()
		if err == nil {
			config.DB.Model(&reg).Update("calendar_event_id", created.Id)
		}
	}
}
