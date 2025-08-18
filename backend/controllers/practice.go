package controllers

import (
	"net/http"
	"time"
	"cfalarm/config"
	"cfalarm/models"
	"github.com/gin-gonic/gin"
)

func GetTodayProblems(c *gin.Context) {
	userID := c.GetUint("user_id")
	var problems []models.PracticeProblem
	today := time.Now().Format("2006-01-02")
	config.DB.Where("user_id = ? AND date_assigned = ?", userID, today).
		Find(&problems)
	c.JSON(http.StatusOK, problems)
}

func GetRecentProblems(c *gin.Context) {
	userID := c.GetUint("user_id")
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
