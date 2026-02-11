package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/Gatete-Bruno/besend/pkg/database"
	"github.com/gin-gonic/gin"
	emailv1alpha1 "github.com/Gatete-Bruno/besend/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	k8sClient client.Client
	k8sOnce   sync.Once
	k8sErr    error
)

func getK8sClient() (client.Client, error) {
	k8sOnce.Do(func() {
		cfg, err := config.GetConfig()
		if err != nil {
			k8sErr = err
			return
		}

		scheme.AddToScheme(scheme.Scheme)
		emailv1alpha1.AddToScheme(scheme.Scheme)

		k8sClient, k8sErr = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	})
	return k8sClient, k8sErr
}

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

	// Get k8s client lazily
	k8sClient, err := getK8sClient()
	if err != nil {
		database.UpdateEmailStatus(email.ID, "failed", stringPtr(fmt.Sprintf("k8s unavailable: %v", err)))
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Email service temporarily unavailable"})
		return
	}

	// Create EmailSenderConfig in Kubernetes if it doesn't exist
	configName := fmt.Sprintf("config-%d-%d", customer.ID, req.SMTPConfigID)
	emailSenderConfig := &emailv1alpha1.EmailSenderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configName,
			Namespace: "besend",
		},
		Spec: emailv1alpha1.EmailSenderConfigSpec{
			Provider:          "native-smtp",
			SenderEmail:       smtpConfig.FromEmail,
			Domain:            smtpConfig.SMTPHost,
			Port:              smtpConfig.SMTPPort,
			APITokenSecretRef: "mailhog-smtp",
		},
	}

	ctx := context.Background()
	if err := k8sClient.Create(ctx, emailSenderConfig); err != nil {
		// Config might already exist, that's fine
	}

	// Create Email CRD for the operator to process
	emailCRD := &emailv1alpha1.Email{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("email-%d", email.ID),
			Namespace: "besend",
		},
		Spec: emailv1alpha1.EmailSpec{
			SenderConfigRef: configName,
			RecipientEmail:  req.To,
			Subject:         req.Subject,
			Body:            req.Body,
		},
	}

	if err := k8sClient.Create(ctx, emailCRD); err != nil {
		database.UpdateEmailStatus(email.ID, "failed", stringPtr(fmt.Sprintf("failed to create k8s resource: %v", err)))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue email"})
		return
	}

	// Increment quota
	if err := customer.IncrementEmailCount(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quota"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":     "Email queued for delivery",
		"email_id":    email.ID,
		"k8s_resource": fmt.Sprintf("email-%d", email.ID),
		"status":      "pending",
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

func stringPtr(s string) *string {
	return &s
}
