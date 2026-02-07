package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/Gatete-Bruno/besend/pkg/database"
    "github.com/Gatete-Bruno/besend/pkg/handlers"
    "github.com/Gatete-Bruno/besend/pkg/middleware"
)

func main() {
    // Database connection
    dbHost := getEnv("DATABASE_HOST", "localhost")
    dbPort := getEnv("DATABASE_PORT", "5432")
    dbUser := getEnv("DATABASE_USER", "besend")
    dbPassword := getEnv("DATABASE_PASSWORD", "besend")
    dbName := getEnv("DATABASE_NAME", "besend")

    db, err := database.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Run migrations
    if err := database.Migrate(db); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Initialize router
    r := gin.Default()

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })

    // API v1 routes
    v1 := r.Group("/api/v1")
    {
        // Public routes
        customers := v1.Group("/customers")
        {
            customers.POST("", handlers.CreateCustomer(db))
            customers.GET("/:id", handlers.GetCustomer(db))
        }

        // Protected routes (require API key)
        protected := v1.Group("")
        protected.Use(middleware.AuthMiddleware(db))
        {
            // SMTP Configuration
            protected.POST("/smtp", handlers.CreateSMTPConfig(db))
            protected.GET("/smtp", handlers.GetSMTPConfigs(db))
            protected.PUT("/smtp/:id", handlers.UpdateSMTPConfig(db))
            protected.DELETE("/smtp/:id", handlers.DeleteSMTPConfig(db))

            // Email Sending
            protected.POST("/emails/send", handlers.SendEmail(db))
            protected.GET("/emails", handlers.GetEmails(db))
            protected.GET("/emails/:id", handlers.GetEmailByID(db))

            // Customer Info
            protected.GET("/me", handlers.GetCurrentCustomer(db))
            protected.GET("/usage", handlers.GetUsage(db))
        }
    }

    // Start server
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
