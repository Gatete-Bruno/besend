package handlers

import (
        "context"
        "crypto/tls"
        "fmt"
        "net"
        "net/smtp"
        "net/http"
        "strconv"
        "time"

        "besend/pkg/database"

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

        ctx := context.Background()
        timeout := 30 * time.Second
        addr := fmt.Sprintf("%s:%d", smtpConfig.SMTPHost, smtpConfig.SMTPPort)

        conn, err := net.DialTimeout("tcp", addr, timeout)
        if err != nil {
                errorMsg := fmt.Sprintf("dial failed: %v", err)
                database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
                c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
                return
        }
        defer conn.Close()

        client, err := smtp.NewClient(conn, smtpConfig.SMTPHost)
        if err != nil {
                errorMsg := fmt.Sprintf("smtp client failed: %v", err)
                database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
                c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
                return
        }
        defer client.Close()

        if ok, _ := client.Extension("STARTTLS"); ok {
                tlsConfig := &tls.Config{ServerName: smtpConfig.SMTPHost}
                if err := client.StartTLS(tlsConfig); err != nil {
                        errorMsg := fmt.Sprintf("starttls failed: %v", err)
                        database.UpdateEmailStatus(email.ID, "failed", &errorMsg)
                        c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
                        return
                }
        }

        if smtpConfig.Username != "" {
                auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.SMTPHost)
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

        if smtpConfig.SMTPHost != "localhost" && smtpConfig.SMTPHost != "127.0.0.1" {
                _ = client.Quit()
        }

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
