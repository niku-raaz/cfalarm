package controllers

import (
	"cfalarm/config"
	"cfalarm/models"
	"cfalarm/services"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// STEP 1: Redirect user to Google
func GoogleLogin(c *gin.Context) {
	// Use the function to get a fresh config with loaded env vars
	googleOAuthConfig := services.GetGoogleOAuthConfig()
	url := googleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// STEP 2: Google redirects back here with ?code=XXXX
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	// Also use the function here
	googleOAuthConfig := services.GetGoogleOAuthConfig()

	// Exchange code for token
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code"})
		return
	}

	tokenJSON, _ := json.Marshal(token)

	// Fetch user info from Google API
	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	var gUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Check if user exists, otherwise create
	var user models.User
	db := config.DB
	if err := db.Where("email = ?", gUser.Email).First(&user).Error; err != nil {
		// create new user
		user = models.User{
			Name:        gUser.Name,
			Email:       gUser.Email,
			GoogleToken: string(tokenJSON),
		}
		db.Create(&user)
	} else {
		db.Model(&user).Update("google_token", string(tokenJSON))
	}

	// Generate JWT (valid for 7 days)
	exp := time.Now().Add(7 * 24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     exp.Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := jwtToken.SignedString([]byte(config.GetEnv("JWT_SECRET")))

	// Return token to frontend (or set cookie if you prefer)
	//c.JSON(http.StatusOK, gin.H{"token": tokenString})
	frontendURL := "http://localhost:3000/auth/callback?token=" + tokenString
    c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}