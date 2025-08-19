package controllers

import (
	"cfalarm/services"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
)

func GetUpcomingContests(c *gin.Context) {
	// Call the new service to fetch contests
	contests, err := services.FetchUpcomingContests()
	if err != nil {
		// If there's an error, return an internal server error response
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contests"})
		return
	}
	// Return the list of upcoming contests as JSON
	log.Printf("Number of upcoming contests: %d", len(contests))
	c.JSON(http.StatusOK, contests)
}
func GetFinishedContests(c *gin.Context) {
	contests, err := services.FetchFinishedContests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch finished contests"})
		return
	}
	c.JSON(http.StatusOK, contests)
}

func RegisterForContest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented"})
}

// Cron
func AutoRegisterCron(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "cron placeholder"})
}
