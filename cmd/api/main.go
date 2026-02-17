package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Gatete-Bruno/besend/pkg/database"
	"github.com/Gatete-Bruno/besend/pkg/api/handlers"
	"github.com/Gatete-Bruno/besend/pkg/api/middleware"
)

func main() {
	dbHost := getEnv("DATABASE_HOST", "localhost")
	dbPortStr := getEnv("DATABASE_PORT", "5432")
	dbPort, _ := strconv.Atoi(dbPortStr)
	dbUser := getEnv("DATABASE_USER", "besenduser")
	dbPassword := getEnv("DATABASE_PASSWORD", "besend")
	dbName := getEnv("DATABASE_NAME", "besend")
	sslMode := getEnv("SSL_MODE", "disable")

	dbConfig := database.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  sslMode,
	}

	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", handlers.GetCustomerInfo)

			protected.POST("/smtp", handlers.CreateSMTPConfig)
			protected.GET("/smtp", handlers.GetSMTPConfigs)
			protected.DELETE("/smtp/:id", handlers.DeleteSMTPConfig)

			protected.POST("/emails/send", handlers.SendEmail)
			protected.GET("/emails", handlers.GetEmailHistory)
			protected.GET("/emails/stats", handlers.GetEmailStats)

			protected.POST("/keys", handlers.CreateAPIKey)
			protected.GET("/keys", handlers.GetAPIKeys)
			protected.DELETE("/keys/:id", handlers.DeleteAPIKey)
		}
	}

	port := getEnv("API_PORT", "8080")
	log.Printf("API Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
