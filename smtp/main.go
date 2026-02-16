package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Organization struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
}

type Domain struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Domain         string `json:"domain"`
	Verified       bool   `json:"verified"`
	SPF            string `json:"spf_record"`
	DKIM           string `json:"dkim_record"`
	DMARC          string `json:"dmarc_record"`
	CreatedAt      string `json:"created_at"`
}

type SMTPCredential struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Username       string `json:"username"`
	Password       string `json:"-"`
	CreatedAt      string `json:"created_at"`
}

var db *sql.DB

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://besenduser:besend_secure_password_123@postgres.besend.svc.cluster.local:5432/besend?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	createTables()
}

func createTables() {
	schema := `
	CREATE TABLE IF NOT EXISTS organizations (
		id VARCHAR(36) PRIMARY KEY,
		customer_id VARCHAR(36) NOT NULL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS sending_domains (
		id VARCHAR(36) PRIMARY KEY,
		organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id),
		domain VARCHAR(255) NOT NULL UNIQUE,
		verified BOOLEAN DEFAULT FALSE,
		spf_record TEXT,
		dkim_record TEXT,
		dmarc_record TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS smtp_credentials (
		id VARCHAR(36) PRIMARY KEY,
		organization_id VARCHAR(36) NOT NULL REFERENCES organizations(id),
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_org_customer ON organizations(customer_id);
	CREATE INDEX IF NOT EXISTS idx_domain_org ON sending_domains(organization_id);
	CREATE INDEX IF NOT EXISTS idx_creds_org ON smtp_credentials(organization_id);
	`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	log.Println("Database tables initialized")
}

func generateSecurePassword() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func createOrganization(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
		Name       string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := generateID()
	createdAt := time.Now().Format(time.RFC3339)

	_, err := db.Exec(
		"INSERT INTO organizations (id, customer_id, name, created_at) VALUES ($1, $2, $3, $4)",
		id, req.CustomerID, req.Name, createdAt,
	)

	if err != nil {
		log.Printf("Error creating organization: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, Organization{
		ID:         id,
		CustomerID: req.CustomerID,
		Name:       req.Name,
		CreatedAt:  createdAt,
	})
}

func listOrganizations(c *gin.Context) {
	customerID := c.Param("customerID")

	rows, err := db.Query(
		"SELECT id, customer_id, name, created_at FROM organizations WHERE customer_id = $1",
		customerID,
	)
	if err != nil {
		log.Printf("Error querying organizations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}
	defer rows.Close()

	var orgs []Organization
	for rows.Next() {
		var org Organization
		if err := rows.Scan(&org.ID, &org.CustomerID, &org.Name, &org.CreatedAt); err != nil {
			log.Printf("Error scanning organization: %v", err)
			continue
		}
		orgs = append(orgs, org)
	}

	c.JSON(http.StatusOK, orgs)
}

func addDomain(c *gin.Context) {
	orgID := c.Param("orgID")

	var req struct {
		Domain string `json:"domain" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := generateID()
	createdAt := time.Now().Format(time.RFC3339)

	_, err := db.Exec(
		"INSERT INTO sending_domains (id, organization_id, domain, created_at) VALUES ($1, $2, $3, $4)",
		id, orgID, req.Domain, createdAt,
	)

	if err != nil {
		log.Printf("Error adding domain: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add domain"})
		return
	}

	spf := fmt.Sprintf("v=spf1 ip4:YOUR_IP include:besend.rw ~all")
	dkim := "TODO: Generate DKIM public key"
	dmarc := fmt.Sprintf("v=DMARC1; p=quarantine; rua=mailto:admin@%s", req.Domain)

	c.JSON(http.StatusCreated, Domain{
		ID:             id,
		OrganizationID: orgID,
		Domain:         req.Domain,
		Verified:       false,
		SPF:            spf,
		DKIM:           dkim,
		DMARC:          dmarc,
		CreatedAt:      createdAt,
	})
}

func createSMTPCredential(c *gin.Context) {
	orgID := c.Param("orgID")

	var req struct {
		Description string `json:"description"`
	}

	if err := c.BindJSON(&req); err != nil {
		req.Description = "Default credential"
	}

	username := generateID()
	password, _ := generateSecurePassword()

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create credential"})
		return
	}

	id := generateID()
	createdAt := time.Now().Format(time.RFC3339)

	_, err = db.Exec(
		"INSERT INTO smtp_credentials (id, organization_id, username, password, created_at) VALUES ($1, $2, $3, $4, $5)",
		id, orgID, username, hashedPassword, createdAt,
	)

	if err != nil {
		log.Printf("Error creating credential: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create credential"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         id,
		"username":   username,
		"password":   password,
		"created_at": createdAt,
		"note":       "Save this password securely. You cannot retrieve it later.",
	})
}

func listCredentials(c *gin.Context) {
	orgID := c.Param("orgID")

	rows, err := db.Query(
		"SELECT id, organization_id, username, created_at FROM smtp_credentials WHERE organization_id = $1",
		orgID,
	)
	if err != nil {
		log.Printf("Error querying credentials: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credentials"})
		return
	}
	defer rows.Close()

	var creds []SMTPCredential
	for rows.Next() {
		var cred SMTPCredential
		if err := rows.Scan(&cred.ID, &cred.OrganizationID, &cred.Username, &cred.CreatedAt); err != nil {
			log.Printf("Error scanning credential: %v", err)
			continue
		}
		creds = append(creds, cred)
	}

	c.JSON(http.StatusOK, creds)
}

func getAuditLogs(c *gin.Context) {
	customerID := c.Param("customerID")
	limit := c.DefaultQuery("limit", "100")
	offset := c.DefaultQuery("offset", "0")

	rows, err := db.Query(
		`SELECT id, customer_id, postal_message_id, sender_email, recipient_email, subject, status, created_at, delivered_at FROM email_audit_logs WHERE customer_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		customerID, limit, offset,
	)
	if err != nil {
		log.Printf("Error querying audit logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var (
			id             string
			customerId     string
			messageId      string
			senderEmail    sql.NullString
			recipientEmail sql.NullString
			subject        sql.NullString
			status         string
			createdAt      string
			deliveredAt    sql.NullString
		)

		if err := rows.Scan(&id, &customerId, &messageId, &senderEmail, &recipientEmail, &subject, &status, &createdAt, &deliveredAt); err != nil {
			log.Printf("Error scanning audit log: %v", err)
			continue
		}

		log := map[string]interface{}{
			"id":           id,
			"customer_id":  customerId,
			"message_id":   messageId,
			"sender":       senderEmail.String,
			"recipient":    recipientEmail.String,
			"subject":      subject.String,
			"status":       status,
			"created_at":   createdAt,
			"delivered_at": deliveredAt.String,
		}
		logs = append(logs, log)
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"limit":  limit,
		"offset": offset,
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func main() {
	router := gin.Default()

	router.GET("/health", healthCheck)
	router.POST("/api/organizations", createOrganization)
	router.GET("/api/organizations/:customerID", listOrganizations)
	router.POST("/api/organizations/:orgID/domains", addDomain)
	router.POST("/api/organizations/:orgID/credentials", createSMTPCredential)
	router.GET("/api/organizations/:orgID/credentials", listCredentials)
	router.GET("/api/audit-logs/:customerID", getAuditLogs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting SMTP management service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
