package controllers

import (
	"cfalarm/config"
	"cfalarm/models"
	"cfalarm/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPracticeProblems(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Get rating from query param, default to 1200 if not provided
	ratingStr := c.DefaultQuery("rating", "1200")
	rating, _ := strconv.Atoi(ratingStr)

	// 1. Get user from DB to find their Codeforces handle
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.CodeforcesID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please set your Codeforces ID in your profile."})
		return
	}

	// 2. Fetch practice problems from the service
	problems, err := services.GetPracticeProblems(user.CodeforcesID, rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch practice problems"})
		return
	}

	c.JSON(http.StatusOK, problems)
}

func GetTodayProblems(c *gin.Context) {
	userID := c.GetUint("user_id")
	var problems []models.PracticeProblem
	today := time.Now().Format("2006-01-02")
	config.DB.Where("user_id = ? AND date_assigned = ?", userID, today).
		Find(&problems)
	c.JSON(http.StatusOK, problems)
}

func GetRecentProblems(c *gin.Context) {
	//userID := c.GetUint("user_id")
	var problems []models.PracticeProblem
	config.DB.Where("user_id = ? AND status = 'unsolved'").
		Order("date_assigned desc").Limit(10).Find(&problems)
	c.JSON(http.StatusOK, problems)
}

func MarkSolved(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct{ ProblemID string }
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid"})
		return
	}
	now := time.Now()
	config.DB.Model(&models.PracticeProblem{}).
		Where("user_id = ? AND problem_id = ?", userID, req.ProblemID).
		Updates(models.PracticeProblem{Status: "solved", SolvedAt: &now})

	c.JSON(http.StatusOK, gin.H{"message": "marked"})
}

// TODO: FetchDailyProblems should be implemented in cron/service
func FetchDailyProblems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "cron placeholder"})
}
