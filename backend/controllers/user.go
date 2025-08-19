package controllers

import (
	"net/http"
	"cfalarm/config"
	"cfalarm/models"
	"cfalarm/services"
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

func VerifyAndSaveCFKeys(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		CodeforcesID string `json:"codeforcesId"`
		APIKey       string `json:"apiKey"`
		APISecret    string `json:"apiSecret"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Step 1: Verify the handle and keys with the Codeforces API
	isValid, err := services.VerifyCFKeys(req.CodeforcesID, req.APIKey, req.APISecret)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials. Please check your Handle, API Key, and Secret."})
		return
	}

	// Step 2: If valid, save all details to the database
	userUpdates := models.User{
		CodeforcesID: req.CodeforcesID,
		CfApiKey:     req.APIKey,
		CfApiSecret:  req.APISecret, // Note: For production, you should encrypt this secret before saving
	}
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(userUpdates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Codeforces credentials verified and saved successfully!"})
}

func VerifyAndSaveCFProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		CodeforcesID string `json:"codeforcesId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Step 1: Verify the handle with the Codeforces API
	isValid, err := services.VerifyCodeforcesHandle(req.CodeforcesID)
	if err != nil || !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-existent Codeforces handle."})
		return
	}

	// Step 2: If valid, save it to the database
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("codeforces_id", req.CodeforcesID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Codeforces ID verified and saved successfully!"})
}

func UpdateProfile(c *gin.Context) {
    userID := c.GetUint("user_id")
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
