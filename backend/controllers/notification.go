package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func SendReminderEmail(c *gin.Context) {
	// will use email_service in /services
	c.JSON(http.StatusOK, gin.H{"message": "cron placeholder"})
}
