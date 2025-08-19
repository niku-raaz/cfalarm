package main

import (
	"cfalarm/config"
	"cfalarm/controllers"
	"cfalarm/cron"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	_, err := config.SetupDatabase()
	if err != nil {
		panic("Failed to connect to database")
	}

	//db.AutoMigrate(&models.User{})

	router := gin.Default()
	router.Use(config.CORSMiddleware())

	api := router.Group("/api")
	{
		api.GET("/auth/google", controllers.GoogleLogin)
		api.GET("/auth/google/callback", controllers.GoogleCallback)
		api.GET("/user/me", controllers.AuthMiddleware(), controllers.GetProfile)
		api.PUT("/user/me", controllers.AuthMiddleware(), controllers.UpdateProfile)

		api.POST("/user/verify-cf-handle", controllers.AuthMiddleware(), controllers.VerifyAndSaveCFHandle)

		api.GET("/practice/problems", controllers.AuthMiddleware(), controllers.GetPracticeProblems)
		api.GET("/practice/today", controllers.AuthMiddleware(), controllers.GetTodayProblems)

		api.GET("/practice/recent", controllers.AuthMiddleware(), controllers.GetRecentProblems)
		api.POST("/practice/mark-solved", controllers.AuthMiddleware(), controllers.MarkSolved)
		api.POST("/practice/fetch-daily", controllers.FetchDailyProblems)

		api.GET("/contests/upcoming", controllers.GetUpcomingContests)
		api.GET("/contests/finished", controllers.GetFinishedContests)

		api.POST("/notifications/send-reminder", controllers.SendReminderEmail)
	}

	// Start cron jobs
	cron.Start()

	router.Run(":" + config.GetEnv("PORT"))
}
