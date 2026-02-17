package database

type Customer struct {
	ID                  int
	Email               string
	PasswordHash        string
	Plan                string
	MonthlyQuota        int
	EmailsSentThisMonth int
	Active              bool
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
}

type Email struct {
	ID           int
	CustomerID   int
	SMTPConfigID int
	ToEmail      string
	Subject      string
	Body         string
	Status       string
}

type APIKey struct {
	ID        int
	Name      string
	CreatedAt string
	LastUsedAt *string
}

func CreateCustomer(email, plan string) (*Customer, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO customers (email, plan) VALUES ($1, $2) RETURNING id",
		email, plan,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &Customer{ID: id, Email: email, Plan: plan}, nil
}

func GetCustomer(id int) (*Customer, error) {
	var c Customer
	err := DB.QueryRow(
		"SELECT id, email, plan, monthly_quota, emails_sent_this_month, active FROM customers WHERE id = $1",
		id,
	).Scan(&c.ID, &c.Email, &c.Plan, &c.MonthlyQuota, &c.EmailsSentThisMonth, &c.Active)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func CreateSMTPConfig(customerID int, name, host string, port int, username, password, fromEmail string) (*SMTPConfig, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO smtp_configs (customer_id, name, smtp_host, smtp_port, username, password, from_email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		customerID, name, host, port, username, password, fromEmail,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &SMTPConfig{ID: id, CustomerID: customerID, Name: name, SMTPHost: host, SMTPPort: port, Username: username, Password: password, FromEmail: fromEmail}, nil
}

func GetSMTPConfig(customerID int) (*SMTPConfig, error) {
	var s SMTPConfig
	err := DB.QueryRow(
		"SELECT id, customer_id, name, smtp_host, smtp_port, username, password, from_email FROM smtp_configs WHERE customer_id = $1 LIMIT 1",
		customerID,
	).Scan(&s.ID, &s.CustomerID, &s.Name, &s.SMTPHost, &s.SMTPPort, &s.Username, &s.Password, &s.FromEmail)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetSMTPConfigs(customerID int) ([]SMTPConfig, error) {
	rows, err := DB.Query(
		"SELECT id, customer_id, name, smtp_host, smtp_port, username, password, from_email FROM smtp_configs WHERE customer_id = $1",
		customerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []SMTPConfig
	for rows.Next() {
		var s SMTPConfig
		if err := rows.Scan(&s.ID, &s.CustomerID, &s.Name, &s.SMTPHost, &s.SMTPPort, &s.Username, &s.Password, &s.FromEmail); err != nil {
			return nil, err
		}
		configs = append(configs, s)
	}
	return configs, nil
}

func CreateEmail(customerID, smtpConfigID int, toEmail, subject, body string) (*Email, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO emails (customer_id, smtp_config_id, to_email, subject, body) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		customerID, smtpConfigID, toEmail, subject, body,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &Email{ID: id, CustomerID: customerID, SMTPConfigID: smtpConfigID, ToEmail: toEmail, Subject: subject, Body: body, Status: "pending"}, nil
}

func UpdateEmailStatus(emailID int, status string, errorMsg *string) error {
	if errorMsg != nil {
		_, err := DB.Exec(
			"UPDATE emails SET status = $1, error_message = $2 WHERE id = $3",
			status, *errorMsg, emailID,
		)
		return err
	}
	_, err := DB.Exec(
		"UPDATE emails SET status = $1, sent_at = CURRENT_TIMESTAMP WHERE id = $2",
		status, emailID,
	)
	return err
}

func GetEmailsByCustomer(customerID, limit, offset int) ([]Email, error) {
	rows, err := DB.Query(
		"SELECT id, customer_id, smtp_config_id, to_email, subject, body, status FROM emails WHERE customer_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		customerID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var e Email
		if err := rows.Scan(&e.ID, &e.CustomerID, &e.SMTPConfigID, &e.ToEmail, &e.Subject, &e.Body, &e.Status); err != nil {
			return nil, err
		}
		emails = append(emails, e)
	}
	return emails, nil
}

func GetEmailStats(customerID int) (map[string]interface{}, error) {
	var sent, pending, failed int64
	err := DB.QueryRow(
		"SELECT COUNT(CASE WHEN status = 'sent' THEN 1 END), COUNT(CASE WHEN status = 'pending' THEN 1 END), COUNT(CASE WHEN status = 'failed' THEN 1 END) FROM emails WHERE customer_id = $1",
		customerID,
	).Scan(&sent, &pending, &failed)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"sent":    sent,
		"pending": pending,
		"failed":  failed,
	}, nil
}

func (c *Customer) IncrementEmailCount() error {
	_, err := DB.Exec(
		"UPDATE customers SET emails_sent_this_month = emails_sent_this_month + 1 WHERE id = $1",
		c.ID,
	)
	return err
}

func GetCustomerByEmail(email string) (*Customer, error) {
	var c Customer
	err := DB.QueryRow(
		"SELECT id, email, password_hash, plan, monthly_quota, emails_sent_this_month, active FROM customers WHERE email = $1",
		email,
	).Scan(&c.ID, &c.Email, &c.PasswordHash, &c.Plan, &c.MonthlyQuota, &c.EmailsSentThisMonth, &c.Active)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetCustomerByAPIKey(keyHash string) (*Customer, error) {
	var c Customer
	err := DB.QueryRow(
		"SELECT c.id, c.email, c.password_hash, c.plan, c.monthly_quota, c.emails_sent_this_month, c.active FROM customers c JOIN api_keys ak ON c.id = ak.customer_id WHERE ak.key_hash = $1",
		keyHash,
	).Scan(&c.ID, &c.Email, &c.PasswordHash, &c.Plan, &c.MonthlyQuota, &c.EmailsSentThisMonth, &c.Active)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func CreateAPIKey(customerID int, keyHash, name string) (*APIKey, error) {
	var id int
	var createdAt string
	err := DB.QueryRow(
		"INSERT INTO api_keys (customer_id, key_hash, name) VALUES ($1, $2, $3) RETURNING id, created_at",
		customerID, keyHash, name,
	).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	return &APIKey{ID: id, Name: name, CreatedAt: createdAt}, nil
}

func GetAPIKeysByCustomer(customerID int) ([]APIKey, error) {
	rows, err := DB.Query(
		"SELECT id, name, created_at, last_used_at FROM api_keys WHERE customer_id = $1",
		customerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var k APIKey
		if err := rows.Scan(&k.ID, &k.Name, &k.CreatedAt, &k.LastUsedAt); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

func DeleteAPIKey(keyID int) error {
	_, err := DB.Exec("DELETE FROM api_keys WHERE id = $1", keyID)
	return err
}

func DeleteSMTPConfig(configID string, customerID int) error {
	_, err := DB.Exec(
		"DELETE FROM smtp_configs WHERE id = $1 AND customer_id = $2",
		configID, customerID,
	)
	return err
}
