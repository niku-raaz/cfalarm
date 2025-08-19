package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cfalarm/config"
	"cfalarm/models"
	"cfalarm/services"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func RunAutoContestRegistration() {
	log.Println("[cron] Starting auto calendar sync process...")

	upcomingContests, err := services.FetchUpcomingContests()
	if err != nil {
		log.Println("Error fetching upcoming contests:", err)
		return
	}
	if len(upcomingContests) == 0 {
		log.Println("No upcoming contests to sync.")
		return
	}

	var users []models.User
	config.DB.Find(&users) // Get all users

	if len(users) == 0 {
		log.Println("No users found for calendar sync.")
		return
	}

	for _, user := range users {
		log.Printf("Processing user: %s", user.Email)
		for _, contest := range upcomingContests {
			// Directly add the contest to the user's calendar
			addContestToCalendar(user, contest)
		}
	}
	log.Println("[cron] Auto calendar sync process finished.")
}

// Helper function to add the event to Google Calendar
func addContestToCalendar(user models.User, contest services.CFContest) {
	if user.GoogleToken == "" {
		log.Printf("User %s has no Google token, skipping calendar event.", user.Email)
		return
	}

	var token oauth2.Token
	if err := json.Unmarshal([]byte(user.GoogleToken), &token); err != nil {
		log.Printf("Error unmarshaling Google token for user %s: %v", user.Email, err)
		return
	}

	ctx := context.Background()
	config := services.GetGoogleOAuthConfig()
	calendarService, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, &token)))
	if err != nil {
		log.Printf("Error creating calendar service for %s: %v", user.Email, err)
		return
	}

	startTime := time.Unix(contest.StartTimeSeconds, 0)
	endTime := startTime.Add(time.Duration(contest.DurationSeconds) * time.Second)

	event := &calendar.Event{
		Summary:     contest.Name,
		Description: fmt.Sprintf("Codeforces contest. URL: https://codeforces.com/contest/%d", contest.ID),
		Start:       &calendar.EventDateTime{DateTime: startTime.Format(time.RFC3339)},
		End:         &calendar.EventDateTime{DateTime: endTime.Format(time.RFC3339)},
	}

	_, err = calendarService.Events.Insert("primary", event).Do()
	if err != nil {
		log.Printf("Could not create calendar event for %s: %v", user.Email, err)
	} else {
		log.Printf("Successfully added contest '%s' to Google Calendar for %s", contest.Name, user.Email)
	}
}