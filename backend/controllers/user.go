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
	c.JSON(http.StatusOK, gin.H{
        "codeforces_id": u.CodeforcesID,
        "codechef_id":   u.CodechefID,
    })
}

func UpdateProfile(c *gin.Context) {
    userID := c.GetInt("user_id")
    var req struct {
        CodeforcesID string `json:"codeforcesId"`
        CodechefID   string `json:"codechefId"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    config.DB.Model(&models.User{}).Where("id = ?", userID).
        Updates(models.User{
            CodeforcesID: req.CodeforcesID,
            CodechefID:   req.CodechefID,
        })
    c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}
