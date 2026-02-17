package handlers

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"net/http"
	"strconv"
	"time"

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

	configs, err := database.GetSMTPConfigs(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch configs"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func SendEmail(c *gin.Context) {
	var req struct {
		To       string `json:"to" binding:"required"`
		Subject  string `json:"subject" binding:"required"`
		Body     string `json:"body" binding:"required"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := c.MustGet("customer").(*database.Customer)

	if customer.EmailsSentThisMonth >= customer.MonthlyQuota {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Monthly quota exceeded",
			"quota": customer.MonthlyQuota,
			"used":  customer.EmailsSentThisMonth,
		})
		return
	}

	smtpConfig, err := database.GetSMTPConfig(customer.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SMTP configuration not found"})
		return
	}

	email, err := database.CreateEmail(customer.ID, smtpConfig.ID, req.To, req.Subject, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	harakaHost := "haraka-smtp.smtp.svc.cluster.local"
	harakaPort := 25

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", harakaHost, harakaPort), 5*time.Second)
	if err != nil {
		errorMsg := fmt.Sprintf("connection to Haraka failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	client, err := smtp.NewClient(conn, harakaHost)
	if err != nil {
		errorMsg := fmt.Sprintf("smtp client failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	defer client.Close()

	if req.Username != "" && req.Password != "" {
		auth := smtp.PlainAuth("", req.Username, req.Password, harakaHost)
		if err := client.Auth(auth); err != nil {
			errorMsg := fmt.Sprintf("auth failed: %v", err)
			database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
			return
		}
	}

	if err := client.Mail(smtpConfig.FromEmail); err != nil {
		errorMsg := fmt.Sprintf("mail from failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	if err := client.Rcpt(req.To); err != nil {
		errorMsg := fmt.Sprintf("rcpt to failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	wc, err := client.Data()
	if err != nil {
		errorMsg := fmt.Sprintf("data failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", smtpConfig.FromEmail, req.To, req.Subject, req.Body)
	if _, err := wc.Write([]byte(msg)); err != nil {
		wc.Close()
		errorMsg := fmt.Sprintf("write failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	if err := wc.Close(); err != nil {
		errorMsg := fmt.Sprintf("close failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	client.Quit()

	if err := database.UpdateEmailStatus(email.ID, "sent", nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	if err := customer.IncrementEmailCount(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quota"})
		return
	}

	_ = ctx

	c.JSON(http.StatusOK, gin.H{
		"message":     "Email sent successfully",
		"email_id":    email.ID,
		"smtp_config": smtpConfig.Name,
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

func DeleteSMTPConfig(c *gin.Context) {
	configID := c.Param("id")

	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Config ID required"})
		return
	}

	customer := c.MustGet("customer").(*database.Customer)

	if err := database.DeleteSMTPConfig(configID, customer.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMTP config deleted"})
}
