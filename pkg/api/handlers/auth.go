package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/Gatete-Bruno/besend/pkg/auth"
	"github.com/Gatete-Bruno/besend/pkg/database"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Plan     string `json:"plan"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Plan == "" {
		req.Plan = "starter"
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	customer, err := database.CreateCustomer(req.Email, passwordHash)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	token, err := auth.GenerateToken(customer.ID, customer.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer registered successfully",
		"customer": gin.H{
			"id":    customer.ID,
			"email": customer.Email,
			"plan":  customer.Plan,
		},
		"token": token,
	})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := database.GetCustomerByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !auth.VerifyPassword(customer.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(customer.ID, customer.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"customer": gin.H{
			"id":    customer.ID,
			"email": customer.Email,
			"plan":  customer.Plan,
		},
		"token": token,
	})
}

func CreateAPIKey(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)

	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plainKey := make([]byte, 32)
	if _, err := rand.Read(plainKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate key"})
		return
	}
	keyString := hex.EncodeToString(plainKey)

	keyHash := sha256.Sum256([]byte(keyString))
	keyHashStr := hex.EncodeToString(keyHash[:])

	apiKey, err := database.CreateAPIKey(customer.ID, keyHashStr, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "API key created",
		"key": gin.H{
			"id":         apiKey.ID,
			"name":       apiKey.Name,
			"key":        keyString,
			"created_at": apiKey.CreatedAt,
		},
		"warning": "Save this key safely - you won't be able to see it again",
	})
}

func GetAPIKeys(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)

	keys, err := database.GetAPIKeysByCustomer(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch API keys"})
		return
	}

	c.JSON(http.StatusOK, keys)
}

func DeleteAPIKey(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)
	
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := database.DeleteAPIKey(customer.ID, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted"})
}
