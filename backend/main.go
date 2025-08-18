package main

import (
	"cfalarm/config"
	"cfalarm/controllers"
	"cfalarm/cron"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
    
	clientID := config.GetEnv("GOOGLE_CLIENT_ID")
	log.Printf("DEBUG: Loaded GOOGLE_CLIENT_ID is: [%s]", clientID)
	redirectURL := config.GetEnv("GOOGLE_REDIRECT_URL")
	log.Printf("DEBUG: Loaded GOOGLE_REDIRECT_URL is: [%s]", redirectURL)

	_, err := config.SetupDatabase()
	if err != nil {
		panic("Failed to connect to database")
	}

	router := gin.Default()
	router.Use(config.CORSMiddleware())

	api := router.Group("/api")
	{
		api.GET("/auth/google", controllers.GoogleLogin)
		api.GET("/auth/google/callback", controllers.GoogleCallback)
		api.GET("/user/me", controllers.AuthMiddleware(), controllers.GetProfile)
		api.PUT("/user/me", controllers.AuthMiddleware(), controllers.UpdateProfile)

		//api.GET("/ratings", controllers.AuthMiddleware(), controllers.GetRatings)
		//api.POST("/ratings/fetch", controllers.AuthMiddleware(), controllers.FetchRatings)

		api.GET("/practice/today", controllers.AuthMiddleware(), controllers.GetTodayProblems)
		api.GET("/practice/recent", controllers.AuthMiddleware(), controllers.GetRecentProblems)
		api.POST("/practice/mark-solved", controllers.AuthMiddleware(), controllers.MarkSolved)
		api.POST("/practice/fetch-daily", controllers.FetchDailyProblems)

		api.GET("/contests/upcoming", controllers.AuthMiddleware(), controllers.GetUpcomingContests)
		api.POST("/contests/register", controllers.AuthMiddleware(), controllers.RegisterForContest)
		api.POST("/contests/auto-register", controllers.AutoRegisterCron)

		api.POST("/notifications/send-reminder", controllers.SendReminderEmail)
	}

	// Start cron jobs
	cron.Start()

	router.Run(":" + config.GetEnv("PORT"))
}
