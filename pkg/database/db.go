package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DB *sql.DB

func Connect(cfg Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Println("Connected to PostgreSQL")
	return nil
}

func InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		api_key VARCHAR(64) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		plan VARCHAR(50) DEFAULT 'starter',
		monthly_quota INTEGER DEFAULT 1000,
		emails_sent_this_month INTEGER DEFAULT 0,
		stripe_customer_id VARCHAR(255),
		active BOOLEAN DEFAULT true
	);

	CREATE TABLE IF NOT EXISTS smtp_configs (
		id SERIAL PRIMARY KEY,
		customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		smtp_host VARCHAR(255) NOT NULL,
		smtp_port INTEGER NOT NULL,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		from_email VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(customer_id, name)
	);

	CREATE TABLE IF NOT EXISTS emails (
		id SERIAL PRIMARY KEY,
		customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
		smtp_config_id INTEGER REFERENCES smtp_configs(id) ON DELETE SET NULL,
		to_email VARCHAR(255) NOT NULL,
		subject VARCHAR(500) NOT NULL,
		body TEXT NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		sent_at TIMESTAMP,
		error_message TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_emails_customer_id ON emails(customer_id);
	CREATE INDEX IF NOT EXISTS idx_emails_status ON emails(status);
	CREATE INDEX IF NOT EXISTS idx_smtp_configs_customer_id ON smtp_configs(customer_id);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database schema initialized")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
