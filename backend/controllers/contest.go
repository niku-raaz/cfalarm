package controllers

import (
	"net/http"
	"cfalarm/config"
	"github.com/gin-gonic/gin"
)

func GetUpcomingContests(c *gin.Context) {
	// placeholder until we implement CF API
	c.JSON(http.StatusOK, gin.H{"message": "not implemented"})
}

func RegisterForContest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented"})
}

// Cron
func AutoRegisterCron(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "cron placeholder"})
}
