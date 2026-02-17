package database

import (
	"database/sql"
	"time"
)

type Customer struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	Plan         string
	MonthlyQuota int
}

type SMTPConfig struct {
	ID        int
	CustomerID int
	Name      string
	SMTPHost  string
	SMTPPort  int
	Username  string
	Password  string
	FromEmail string
	CreatedAt time.Time
}

type Email struct {
	ID           int
	CustomerID   int
	SMTPConfigID int
	ToEmail      string
	Subject      string
	Body         string
	Status       string
	CreatedAt    time.Time
	SentAt       *time.Time
	ErrorMessage *string
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

func GetEmailStats(customerID int) (map[string]interface{}, error) {
	var sent, pending, failed int64
	err := DB.QueryRow(`
		SELECT 
			COUNT(CASE WHEN status = 'sent' THEN 1 END),
			COUNT(CASE WHEN status = 'pending' THEN 1 END),
			COUNT(CASE WHEN status = 'failed' THEN 1 END)
		FROM emails
		WHERE customer_id = $1
	`, customerID).Scan(&sent, &pending, &failed)
	
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"sent":    sent,
		"pending": pending,
		"failed":  failed,
	}, nil
}

func DeleteSMTPConfig(customerID, configID int) error {
	_, err := DB.Exec(`
		DELETE FROM smtp_configs
		WHERE id = $1 AND customer_id = $2
	`, configID, customerID)
	return err
}
