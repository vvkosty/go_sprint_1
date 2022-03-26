package app

import (
	"context"
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
		ShortURLID string
		UserID     string
		FullURL    string
		IsDeleted  bool
	}
)

type (
	UniqueViolatesError struct{ Err error }
	EntityDeletedError  struct{ Err error }
)

func (uve *UniqueViolatesError) Error() string {
	return fmt.Sprintf("UniqueViolatesError: %v", uve.Err)
}

func NewUniqueViolatesError(err error) error {
	return &UniqueViolatesError{
		Err: err,
	}
}

func (uve *EntityDeletedError) Error() string {
	return "EntityDeletedError: Row was deleted"
}

func NewPostgresDatabase(dsn string, forceRecreate bool) *PostgresDatabase {
	var md PostgresDatabase
	var err error

	md.db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if forceRecreate {
		md.db.Exec(`DROP TABLE IF EXISTS urls`)
	}

	query := `
		CREATE TABLE IF NOT EXISTS urls(
			short_url_id varchar(50) PRIMARY KEY,
			user_id varchar(50) NOT NULL,
			full_url varchar(255) NOT NULL,
			is_deleted bool NOT NULL DEFAULT false
		)`
	_, err = md.db.Exec(query)
	if err != nil {
		log.Printf("Unable to create tables: %v\n", err)
		os.Exit(1)
	}

	return &md
}

func (m *PostgresDatabase) Find(id string) (string, error) {
	var url URLEntity

	row := m.db.QueryRow("SELECT full_url, is_deleted FROM urls WHERE short_url_id = $1", id)
	err := row.Scan(&url.FullURL, &url.IsDeleted)
	if err != nil {
		return "", err
	}

	if url.IsDeleted {
		return "", &EntityDeletedError{}
	}

	return url.FullURL, nil
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

func (m *PostgresDatabase) DeleteBatchByChecksums(checksums []string) error {
	ctx := context.Background()
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE urls SET is_deleted = true WHERE short_url_id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, checksum := range checksums {
		if _, err = stmt.ExecContext(ctx, checksum); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (m *PostgresDatabase) Close() error {
	return m.db.Close()
}
