package database

import (
	"database/sql"
	"time"
)

type Customer struct {
	ID                  int       `json:"id"`
	Email               string    `json:"email"`
	PasswordHash        string    `json:"-"`
	CreatedAt           time.Time `json:"created_at"`
	Plan                string    `json:"plan"`
	MonthlyQuota        int       `json:"monthly_quota"`
	EmailsSentThisMonth int       `json:"emails_sent_this_month"`
	StripeCustomerID    string    `json:"stripe_customer_id,omitempty"`
	Active              bool      `json:"active"`
}

type APIKey struct {
	ID         int        `json:"id"`
	CustomerID int        `json:"customer_id"`
	KeyHash    string     `json:"-"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
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

func CreateCustomer(email, passwordHash string) (*Customer, error) {
	quota := 1000

	var customer Customer
	err := DB.QueryRow(`
		INSERT INTO customers (email, password_hash, plan, monthly_quota)
		VALUES ($1, $2, 'starter', $3)
		RETURNING id, email, password_hash, created_at, plan, monthly_quota, emails_sent_this_month, active
	`, email, passwordHash, quota).Scan(
		&customer.ID, &customer.Email, &customer.PasswordHash, &customer.CreatedAt,
		&customer.Plan, &customer.MonthlyQuota, &customer.EmailsSentThisMonth, &customer.Active,
	)

	return &customer, err
}

func GetCustomerByEmail(email string) (*Customer, error) {
	var customer Customer
	var stripeID sql.NullString
	err := DB.QueryRow(`
		SELECT id, email, password_hash, created_at, plan, monthly_quota, emails_sent_this_month, stripe_customer_id, active
		FROM customers
		WHERE email = $1 AND active = true
	`, email).Scan(
		&customer.ID, &customer.Email, &customer.PasswordHash, &customer.CreatedAt,
		&customer.Plan, &customer.MonthlyQuota, &customer.EmailsSentThisMonth,
		&stripeID, &customer.Active,
	)
	if stripeID.Valid {
		customer.StripeCustomerID = stripeID.String
	}
	return &customer, err
}

func GetCustomerByAPIKey(keyHash string) (*Customer, error) {
	var customer Customer
	var stripeID sql.NullString
	err := DB.QueryRow(`
		SELECT c.id, c.email, c.password_hash, c.created_at, c.plan, c.monthly_quota, 
		       c.emails_sent_this_month, c.stripe_customer_id, c.active
		FROM customers c
		JOIN api_keys ak ON ak.customer_id = c.id
		WHERE ak.key_hash = $1 AND c.active = true
	`, keyHash).Scan(
		&customer.ID, &customer.Email, &customer.PasswordHash, &customer.CreatedAt,
		&customer.Plan, &customer.MonthlyQuota, &customer.EmailsSentThisMonth,
		&stripeID, &customer.Active,
	)
	if stripeID.Valid {
		customer.StripeCustomerID = stripeID.String
	}

	if err == nil {
		DB.Exec(`UPDATE api_keys SET last_used_at = NOW() WHERE key_hash = $1`, keyHash)
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

func CreateAPIKey(customerID int, keyHash, name string) (*APIKey, error) {
	var apiKey APIKey
	err := DB.QueryRow(`
		INSERT INTO api_keys (customer_id, key_hash, name)
		VALUES ($1, $2, $3)
		RETURNING id, customer_id, key_hash, name, created_at
	`, customerID, keyHash, name).Scan(
		&apiKey.ID, &apiKey.CustomerID, &apiKey.KeyHash, &apiKey.Name, &apiKey.CreatedAt,
	)
	return &apiKey, err
}

func GetAPIKeysByCustomer(customerID int) ([]APIKey, error) {
	rows, err := DB.Query(`
		SELECT id, customer_id, name, created_at, last_used_at
		FROM api_keys
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var key APIKey
		err := rows.Scan(&key.ID, &key.CustomerID, &key.Name, &key.CreatedAt, &key.LastUsedAt)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func DeleteAPIKey(customerID, keyID int) error {
	_, err := DB.Exec(`
		DELETE FROM api_keys
		WHERE id = $1 AND customer_id = $2
	`, keyID, customerID)
	return err
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
