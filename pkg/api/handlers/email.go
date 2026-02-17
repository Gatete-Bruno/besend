package handlers

import (
	"fmt"
	"net"
	"net/smtp"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Gatete-Bruno/besend/pkg/database"
)

func SendEmail(c *gin.Context) {
	var req struct {
		SMTPConfigID int    `json:"smtp_config_id" binding:"required"`
		To           string `json:"to" binding:"required"`
		Subject      string `json:"subject" binding:"required"`
		Body         string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := c.MustGet("customer").(*database.Customer)

	smtpConfig, err := database.GetSMTPConfigByID(req.SMTPConfigID)
	if err != nil || smtpConfig.CustomerID != customer.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "SMTP config not found"})
		return
	}

	email, err := database.CreateEmail(customer.ID, req.SMTPConfigID, req.To, req.Subject, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email"})
		return
	}

	host := "54.77.87.98"
	port := 30587
	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		errorMsg := fmt.Sprintf("connection to Haraka failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		errorMsg := fmt.Sprintf("SMTP client creation failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}
	defer client.Close()

	from := smtpConfig.FromEmail
	if err := client.Mail(from); err != nil {
		errorMsg := fmt.Sprintf("MAIL command failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	if err := client.Rcpt(req.To); err != nil {
		errorMsg := fmt.Sprintf("RCPT command failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	wc, err := client.Data()
	if err != nil {
		errorMsg := fmt.Sprintf("DATA command failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, req.To, req.Subject, req.Body)
	_, err = wc.Write([]byte(msg))
	if err != nil {
		wc.Close()
		errorMsg := fmt.Sprintf("write failed: %v", err)
		database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	err = wc.Close()
	if err != nil {
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

	c.JSON(http.StatusOK, gin.H{
		"message":  "Email sent successfully",
		"email_id": email.ID,
	})
}

func GetEmailHistory(c *gin.Context) {
	customer := c.MustGet("customer").(*database.Customer)
	
	emails, err := database.GetEmailsByCustomer(customer.ID)
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
	customer := c.MustGet("customer").(*database.Customer)
	configID := c.Param("id")
	
	err := database.DeleteSMTPConfig(configID, customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMTP config deleted"})
}
