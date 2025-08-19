package cron

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
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

	var paramList []string
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(paramList)
	paramStringForHash := ""
	for i, p := range paramList {
		if i > 0 {
			paramStringForHash += "&"
		}
		paramStringForHash += p
	}

	randStr := "123456"
	hashInput := fmt.Sprintf("%s/%s?%s#%s", randStr, methodName, paramStringForHash, secret)
	hasher := sha512.New()
	hasher.Write([]byte(hashInput))
	apiSig := hex.EncodeToString(hasher.Sum(nil))

	// --- THIS IS THE FIX ---
	// 1. The base URL does not contain the parameters.
	apiURL := fmt.Sprintf("https://codeforces.com/api/%s", methodName)

	// 2. The parameters are encoded into a form body.
	formData := url.Values{}
	formData.Set("contestId", strconv.Itoa(contestID))
	formData.Set("apiKey", key)
	formData.Set("time", strconv.FormatInt(time.Now().Unix(), 10))
	formData.Set("apiSig", randStr+apiSig) // The signature includes the random prefix

	// 3. Make the POST request with the correct body and content type.
	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	// --- END OF FIX ---

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// If decoding fails, it's likely an HTML error page.
		// We can return a more user-friendly error.
		return fmt.Errorf("registration failed (user may already be registered or registration is closed)")
	}

	if result.Status != "OK" {
		// If it's a valid JSON response but not "OK", use the comment from the API.
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