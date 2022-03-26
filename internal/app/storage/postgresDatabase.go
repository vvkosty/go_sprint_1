package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vvkosty/go_sprint_1/internal/app/helpers"
)

type (
	PostgresDatabase struct {
		db *sql.DB
	}

	URLEntity struct {
		CorrelationID string
		ShortURLID    string
		UserID        string
		FullURL       string
	}
)

type UniqueViolatesError struct{ Err error }

func (uve *UniqueViolatesError) Error() string {
	return fmt.Sprintf("UniqueViolatesError: %v", uve.Err)
}

func NewUniqueViolatesError(err error) error {
	return &UniqueViolatesError{
		Err: err,
	}
}

func NewPostgresDatabase(dsn string) *PostgresDatabase {
	var md PostgresDatabase
	var err error

	md.db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	query := `
		CREATE TABLE IF NOT EXISTS urls(
			short_url_id varchar(50) PRIMARY KEY,
			user_id varchar(50) NOT NULL,
			full_url varchar(255) NOT NULL
		)`
	_, err = md.db.Exec(query)
	if err != nil {
		log.Printf("Unable to create tables: %v\n", err)
		os.Exit(1)
	}

	return &md
}

func (m *PostgresDatabase) Find(id string) (string, error) {
	var url string

	row := m.db.QueryRow("SELECT full_url FROM urls WHERE short_url_id = $1", id)
	err := row.Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (m *PostgresDatabase) Save(url string, userID string) (string, error) {
	checksum := helpers.GenerateChecksum(url)

	_, err := m.db.Exec(
		"INSERT INTO urls (short_url_id, user_id, full_url) VALUES ($1, $2, $3)",
		checksum,
		userID,
		url,
	)

	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return checksum, NewUniqueViolatesError(err)
		}
		return "", err
	}

	return checksum, nil
}

func (m *PostgresDatabase) List(userID string) map[string]string {
	var fullURL string
	var shortURLID string

	result := make(map[string]string)

	rows, err := m.db.Query("SELECT full_url, short_url_id FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return result
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&fullURL, &shortURLID)
		if err != nil {
			return result
		}

		result[shortURLID] = fullURL
	}
	err = rows.Err()
	if err != nil {
		return result
	}

	return result
}

func (m *PostgresDatabase) Close() error {
	return m.db.Close()
}
