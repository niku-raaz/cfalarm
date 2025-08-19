package cron

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"cfalarm/config"
	"cfalarm/models"
	"cfalarm/services"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func RunAutoContestRegistration() {
	log.Println("[cron] Starting auto contest registration process...")

	// 1. Get all upcoming Codeforces contests
	upcomingContests, err := services.FetchUpcomingContests()
	if err != nil {
		log.Println("Error fetching upcoming contests:", err)
		return
	}
	if len(upcomingContests) == 0 {
		log.Println("No upcoming contests to register for.")
		return
	}

	// 2. Find all users who have provided API keys
	var users []models.User
	config.DB.Where("cf_api_key != '' AND cf_api_secret != ''").Find(&users)

	if len(users) == 0 {
		log.Println("No users with API keys to process.")
		return
	}

	// 3. Loop through each user and register them for each contest
	for _, user := range users {
		log.Printf("Processing user: %s", user.Email)
		for _, contest := range upcomingContests {
			// (Optional but recommended) Check if the user is already registered for this contest
			// to avoid unnecessary API calls and potential errors. This would require another
			// API call to contest.standings and checking if the user's handle is in the list.
			// For simplicity, we will attempt to register for all. The API will just return an
			// error if already registered, which we can handle.

			log.Printf("Attempting to register %s for contest '%s'", user.CodeforcesID, contest.Name)
			err := registerUserForContest(user.CodeforcesID, user.CfApiKey, user.CfApiSecret, contest.ID)
			if err != nil {
				// Log the error but continue; it might be that the user is already registered.
				log.Printf("Could not register for contest %d for user %s: %v", contest.ID, user.CodeforcesID, err)
				continue // Move to the next contest
			}

			log.Printf("Successfully registered %s for contest: %s", user.CodeforcesID, contest.Name)

			// 4. Create Google Calendar event
			addContestToCalendar(user, contest)
		}
	}
	log.Println("[cron] Auto contest registration process finished.")
}

// Helper function to perform the registration via Codeforces API
func registerUserForContest(handle, key, secret string, contestID int) error {
	methodName := "contest.register"
	params := map[string]string{
		"contestId": strconv.Itoa(contestID),
		"apiKey":    key,
		"time":      strconv.FormatInt(time.Now().Unix(), 10),
	}

	// Create the signature (apiSig)
	var paramList []string
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(paramList)
	paramString := ""
	for i, p := range paramList {
		if i > 0 {
			paramString += "&"
		}
		paramString += p
	}

	randStr := "123456" // This should be a 6-char random string for production
	hashInput := fmt.Sprintf("%s/%s?%s#%s", randStr, methodName, paramString, secret)
	hasher := sha512.New()
	hasher.Write([]byte(hashInput))
	apiSig := hex.EncodeToString(hasher.Sum(nil))

	finalURL := fmt.Sprintf("https://codeforces.com/api/%s?%s&apiSig=%s%s", methodName, paramString, randStr, apiSig)

	// Codeforces requires a POST request for this method
	resp, err := http.Post(finalURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response from Codeforces: %w", err)
	}

	if result.Status != "OK" {
		return fmt.Errorf("API error: %s", result.Comment)
	}
	return nil
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
	// Use the GetGoogleOAuthConfig function to ensure credentials are loaded
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
		// It's possible the event already exists, which would cause an error.
		// For a more robust system, you could store created event IDs to avoid duplicates.
		log.Printf("Could not create calendar event for %s: %v", user.Email, err)
	} else {
		log.Printf("Successfully added contest '%s' to Google Calendar for %s", contest.Name, user.Email)
	}
}