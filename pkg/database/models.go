package database

import (
        "crypto/rand"
        "database/sql"
        "encoding/hex"
        "time"
)

type Customer struct {
        ID                  int       `json:"id"`
        Email               string    `json:"email"`
        APIKey              string    `json:"api_key"`
        CreatedAt           time.Time `json:"created_at"`
        Plan                string    `json:"plan"`
        MonthlyQuota        int       `json:"monthly_quota"`
        EmailsSentThisMonth int       `json:"emails_sent_this_month"`
        StripeCustomerID    string    `json:"stripe_customer_id,omitempty"`
        Active              bool      `json:"active"`
}

type SMTPConfig struct {
        ID         int       `json:"id"`
        CustomerID int       `json:"customer_id"`
        Name       string    `json:"name"`
        SMTPHost   string    `json:"smtp_host"`
        SMTPPort   int       `json:"smtp_port"`
        Username   string    `json:"username"`
        Password   string    `json:"-"`
        FromEmail  string    `json:"from_email"`
        CreatedAt  time.Time `json:"created_at"`
}

type Email struct {
        ID           int        `json:"id"`
        CustomerID   int        `json:"customer_id"`
        SMTPConfigID *int       `json:"smtp_config_id,omitempty"`
        ToEmail      string     `json:"to_email"`
        Subject      string     `json:"subject"`
        Body         string     `json:"body"`
        Status       string     `json:"status"`
        CreatedAt    time.Time  `json:"created_at"`
        SentAt       *time.Time `json:"sent_at,omitempty"`
        ErrorMessage *string    `json:"error_message,omitempty"`
}

func GenerateAPIKey() (string, error) {
        bytes := make([]byte, 32)
        if _, err := rand.Read(bytes); err != nil {
                return "", err
        }
        return hex.EncodeToString(bytes), nil
}

func CreateCustomer(email, plan string) (*Customer, error) {
        apiKey, err := GenerateAPIKey()
        if err != nil {
                return nil, err
        }

        quota := getQuotaForPlan(plan)

        var customer Customer
        err = DB.QueryRow(`
                INSERT INTO customers (email, api_key, plan, monthly_quota)
                VALUES ($1, $2, $3, $4)
                RETURNING id, email, api_key, created_at, plan, monthly_quota, emails_sent_this_month, active
        `, email, apiKey, plan, quota).Scan(
                &customer.ID, &customer.Email, &customer.APIKey, &customer.CreatedAt,
                &customer.Plan, &customer.MonthlyQuota, &customer.EmailsSentThisMonth, &customer.Active,
        )

        return &customer, err
}

func GetCustomerByAPIKey(apiKey string) (*Customer, error) {
        var customer Customer
        var stripeID sql.NullString
        err := DB.QueryRow(`
                SELECT id, email, api_key, created_at, plan, monthly_quota, emails_sent_this_month, stripe_customer_id, active
                FROM customers
                WHERE api_key = $1 AND active = true
        `, apiKey).Scan(
                &customer.ID, &customer.Email, &customer.APIKey, &customer.CreatedAt,
                &customer.Plan, &customer.MonthlyQuota, &customer.EmailsSentThisMonth,
                &stripeID, &customer.Active,
        )
        if stripeID.Valid {
                customer.StripeCustomerID = stripeID.String
        }
        return &customer, err
}

func (c *Customer) IncrementEmailCount() error {
        _, err := DB.Exec(`
                UPDATE customers
                SET emails_sent_this_month = emails_sent_this_month + 1
                WHERE id = $1
        `, c.ID)
        return err
}

func (c *Customer) HasQuotaAvailable() bool {
        return c.EmailsSentThisMonth < c.MonthlyQuota
}

func CreateSMTPConfig(customerID int, name, host string, port int, username, password, fromEmail string) (*SMTPConfig, error) {
        var config SMTPConfig
        err := DB.QueryRow(`
                INSERT INTO smtp_configs (customer_id, name, smtp_host, smtp_port, username, password, from_email)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                RETURNING id, customer_id, name, smtp_host, smtp_port, username, password, from_email, created_at
        `, customerID, name, host, port, username, password, fromEmail).Scan(
                &config.ID, &config.CustomerID, &config.Name, &config.SMTPHost,
                &config.SMTPPort, &config.Username, &config.Password, &config.FromEmail, &config.CreatedAt,
        )
        return &config, err
}

func GetSMTPConfigsByCustomer(customerID int) ([]SMTPConfig, error) {
        rows, err := DB.Query(`
                SELECT id, customer_id, name, smtp_host, smtp_port, username, from_email, created_at
                FROM smtp_configs
                WHERE customer_id = $1
                ORDER BY created_at DESC
        `, customerID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var configs []SMTPConfig
        for rows.Next() {
                var config SMTPConfig
                err := rows.Scan(
                        &config.ID, &config.CustomerID, &config.Name, &config.SMTPHost,
                        &config.SMTPPort, &config.Username, &config.FromEmail, &config.CreatedAt,
                )
                if err != nil {
                        return nil, err
                }
                configs = append(configs, config)
        }
        return configs, nil
}

func GetSMTPConfigByID(customerID, configID int) (*SMTPConfig, error) {
        var config SMTPConfig
        err := DB.QueryRow(`
                SELECT id, customer_id, name, smtp_host, smtp_port, username, password, from_email, created_at
                FROM smtp_configs
                WHERE id = $1 AND customer_id = $2
        `, configID, customerID).Scan(
                &config.ID, &config.CustomerID, &config.Name, &config.SMTPHost,
                &config.SMTPPort, &config.Username, &config.Password, &config.FromEmail, &config.CreatedAt,
        )
        return &config, err
}

func DeleteSMTPConfig(customerID, configID int) error {
        _, err := DB.Exec(`
                DELETE FROM smtp_configs
                WHERE id = $1 AND customer_id = $2
        `, configID, customerID)
        return err
}

func CreateEmail(customerID int, smtpConfigID *int, toEmail, subject, body string) (*Email, error) {
        var email Email
        err := DB.QueryRow(`
                INSERT INTO emails (customer_id, smtp_config_id, to_email, subject, body, status)
                VALUES ($1, $2, $3, $4, $5, 'pending')
                RETURNING id, customer_id, smtp_config_id, to_email, subject, body, status, created_at
        `, customerID, smtpConfigID, toEmail, subject, body).Scan(
                &email.ID, &email.CustomerID, &email.SMTPConfigID, &email.ToEmail,
                &email.Subject, &email.Body, &email.Status, &email.CreatedAt,
        )
        return &email, err
}

func UpdateEmailStatus(emailID int, status string, errorMsg *string) error {
        sentAt := time.Now()
        _, err := DB.Exec(`
                UPDATE emails
                SET status = $1, sent_at = $2, error_message = $3
                WHERE id = $4
        `, status, sentAt, errorMsg, emailID)
        return err
}

func GetEmailsByCustomer(customerID int, limit, offset int) ([]Email, error) {
        rows, err := DB.Query(`
                SELECT id, customer_id, smtp_config_id, to_email, subject, body, status, created_at, sent_at, error_message
                FROM emails
                WHERE customer_id = $1
                ORDER BY created_at DESC
                LIMIT $2 OFFSET $3
        `, customerID, limit, offset)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var emails []Email
        for rows.Next() {
                var email Email
                err := rows.Scan(
                        &email.ID, &email.CustomerID, &email.SMTPConfigID, &email.ToEmail,
                        &email.Subject, &email.Body, &email.Status, &email.CreatedAt,
                        &email.SentAt, &email.ErrorMessage,
                )
                if err != nil {
                        return nil, err
                }
                emails = append(emails, email)
        }
        return emails, nil
}

func GetEmailStats(customerID int) (map[string]int, error) {
        var total, sent, failed, pending int
        err := DB.QueryRow(`
                SELECT 
                        COUNT(*) as total,
                        COUNT(CASE WHEN status = 'sent' THEN 1 END) as sent,
                        COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
                        COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending
                FROM emails
                WHERE customer_id = $1
        `, customerID).Scan(&total, &sent, &failed, &pending)

        if err != nil {
                return nil, err
        }

        stats := map[string]int{
                "total":   total,
                "sent":    sent,
                "failed":  failed,
                "pending": pending,
        }
        return stats, nil
}

func getQuotaForPlan(plan string) int {
        quotas := map[string]int{
                "starter":      1000,
                "professional": 10000,
                "business":     50000,
                "enterprise":   200000,
        }
        if quota, exists := quotas[plan]; exists {
                return quota
        }
        return 1000
}
