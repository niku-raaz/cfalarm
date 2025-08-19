package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sort" // Import the sort package
	"sync"
	"time"
)

// Represents a single contest from the Codeforces API response
type CFContest struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Phase            string `json:"phase"`
	StartTimeSeconds int64  `json:"startTimeSeconds"`
	DurationSeconds  int    `json:"durationSeconds"`
}

// --- Caching Logic for Upcoming Contests ---
type ContestCache struct {
	Contests []CFContest
	Expiry   time.Time
	mu       sync.Mutex
}
var upcomingContestsCache = &ContestCache{}

// --- Caching Logic for Finished Contests (NEW) ---
var finishedContestsCache = &ContestCache{}
const cacheDuration = 10 * time.Minute


// Fetches UPCOMING contests from the Codeforces API, with caching
func FetchUpcomingContests() ([]CFContest, error) {
	upcomingContestsCache.mu.Lock()
	defer upcomingContestsCache.mu.Unlock()

	if time.Now().Before(upcomingContestsCache.Expiry) {
		log.Println("Returning upcoming contests from cache.")
		return upcomingContestsCache.Contests, nil
	}

	log.Println("Cache expired. Fetching new upcoming contests from Codeforces API.")
	
    // We can call a new helper function to get all contests
	allContests, err := fetchAllContests()
	if err != nil {
		return nil, err
	}

	var upcomingContests []CFContest
	for _, contest := range allContests {
		if contest.Phase == "BEFORE" {
			upcomingContests = append(upcomingContests, contest)
		}
	}

	upcomingContestsCache.Contests = upcomingContests
	upcomingContestsCache.Expiry = time.Now().Add(cacheDuration)

	log.Printf("Number of upcoming contests fetched: %d", len(upcomingContests))
	return upcomingContests, nil
}


// --- NEW FUNCTION to fetch FINISHED contests ---
func FetchFinishedContests() ([]CFContest, error) {
	finishedContestsCache.mu.Lock()
	defer finishedContestsCache.mu.Unlock()

	if time.Now().Before(finishedContestsCache.Expiry) {
		log.Println("Returning finished contests from cache.")
		return finishedContestsCache.Contests, nil
	}

	log.Println("Cache expired. Fetching new finished contests from Codeforces API.")
	allContests, err := fetchAllContests()
	if err != nil {
		return nil, err
	}

	var finishedContests []CFContest
	for _, contest := range allContests {
		if contest.Phase == "FINISHED" {
			finishedContests = append(finishedContests, contest)
		}
	}
    
    // Sort contests by start time in descending order (most recent first)
    sort.Slice(finishedContests, func(i, j int) bool {
        return finishedContests[i].StartTimeSeconds > finishedContests[j].StartTimeSeconds
    })

    // Limit to the last 4 contests
    if len(finishedContests) > 4 {
        finishedContests = finishedContests[:4]
    }

	finishedContestsCache.Contests = finishedContests
	finishedContestsCache.Expiry = time.Now().Add(cacheDuration)

	log.Printf("Number of finished contests fetched: %d", len(finishedContests))
	return finishedContests, nil
}


// --- NEW HELPER FUNCTION to avoid duplicate API calls ---
func fetchAllContests() ([]CFContest, error) {
	url := "https://codeforces.com/api/contest.list"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Status string      `json:"status"`
		Result []CFContest `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Result, nil
}