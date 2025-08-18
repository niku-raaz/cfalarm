package cron

import (
	"log"
	"time"

	"cfalarm/config"
	"cfalarm/models"
	"cfalarm/services"
)

func RunDailyPractice() {
	log.Println("[cron] Fetching daily practice problems...")

	var users []models.User
	config.DB.Find(&users)

	for _, user := range users {
		// TODO: fetch current rating (for now assume a static rating)
		currentRating := 1200

		problems1, _ := services.FetchCFProblemsByRating(currentRating + 100)
		problems2, _ := services.FetchCFProblemsByRating(currentRating + 200)

		// pick first result for simplicity
		if len(problems1) > 0 && len(problems2) > 0 {
			//date := time.Now().Format("2006-01-02")

			pp1 := models.PracticeProblem{
				UserID:       user.ID,
				ProblemID:    problems1[0].Index,
				Rating:       problems1[0].Rating,
				DateAssigned: time.Now(),
				Status:       "unsolved",
			}
			pp2 := models.PracticeProblem{
				UserID:       user.ID,
				ProblemID:    problems2[0].Index,
				Rating:       problems2[0].Rating,
				DateAssigned: time.Now(),
				Status:       "unsolved",
			}

			config.DB.Create(&pp1)
			config.DB.Create(&pp2)
		}
	}
}
