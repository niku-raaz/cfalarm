package controllers

import (
	"net/http"
	"cfalarm/config"
	"cfalarm/models"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var u models.User
	if err := config.DB.First(&u, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	input.ID = userID
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).
		Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}
