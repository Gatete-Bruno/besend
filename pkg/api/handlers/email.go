package handlers

import (
	"net/http"
	"strconv"

	"github.com/Gatete-Bruno/besend/pkg/database"
	"github.com/gin-gonic/gin"
)

func RegisterCustomer(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Plan  string `json:"plan"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Plan == "" {
		req.Plan = "starter"
	}

	customer, err := database.CreateCustomer(req.Email, req.Plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func GetCustomerInfo(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)
	c.JSON(http.StatusOK, customer)
}

func CreateSMTPConfig(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)

	var req struct {
		Name      string `json:"name" binding:"required"`
		SMTPHost  string `json:"smtp_host" binding:"required"`
		SMTPPort  int    `json:"smtp_port" binding:"required"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		FromEmail string `json:"from_email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := database.CreateSMTPConfig(
		customer.ID, req.Name, req.SMTPHost, req.SMTPPort,
		req.Username, req.Password, req.FromEmail,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMTP config"})
		return
	}

	c.JSON(http.StatusCreated, config)
}

func GetSMTPConfigs(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)

	configs, err := database.GetSMTPConfigsByCustomer(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch configs"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func DeleteSMTPConfig(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)
	configID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config ID"})
		return
	}

	if err := database.DeleteSMTPConfig(customer.ID, configID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Config deleted"})
}

func SendEmail(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)

	var req struct {
		SMTPConfigID int    `json:"smtp_config_id" binding:"required"`
		To           string `json:"to" binding:"required,email"`
		Subject      string `json:"subject" binding:"required"`
		Body         string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !customer.HasQuotaAvailable() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Monthly quota exceeded"})
		return
	}

	smtpConfig, err := database.GetSMTPConfigByID(customer.ID, req.SMTPConfigID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SMTP config not found"})
		return
	}

	email, err := database.CreateEmail(customer.ID, &req.SMTPConfigID, req.To, req.Subject, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email"})
		return
	}

	// Queue email for async delivery
	if err := database.UpdateEmailStatus(email.ID, "queued", nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue email"})
		return
	}

	// Increment quota immediately
	if err := customer.IncrementEmailCount(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quota"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Email queued for delivery",
		"email_id":  email.ID,
		"status":    "queued",
	})
}

func GetEmailHistory(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	emails, err := database.GetEmailsByCustomer(customer.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emails"})
		return
	}

	c.JSON(http.StatusOK, emails)
}

func GetEmailStats(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)

	stats, err := database.GetEmailStats(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
