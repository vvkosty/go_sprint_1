package app

import (
	"database/sql"
	"hash/crc32"
	"log"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresDatabase struct {
	db *sql.DB
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
		    correlation_id varchar(50) PRIMARY KEY,
			short_url_id varchar(50) NOT NULL,
			user_id varchar(50) NOT NULL,
			full_url varchar(50) NOT NULL
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

func (m *PostgresDatabase) Save(url string, userId string, correlationId string) (string, error) {
	checksum := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(url))))

	_, err := m.db.Exec(
		"INSERT INTO urls (short_url_id, user_id, full_url, correlation_id) VALUES ($1, $2, $3, $4)",
		checksum,
		userId,
		url,
		correlationId,
	)
	if err != nil {
		return "", err
	}

	return checksum, nil
}

func (m *PostgresDatabase) List(userId string) map[string]string {
	var fullUrl string
	var shortUrlId string

	result := make(map[string]string)

	rows, err := m.db.Query("SELECT full_url, short_url_id FROM urls WHERE user_id = $1", userId)
	if err != nil {
		return result
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&fullUrl, &shortUrlId)
		if err != nil {
			return result
		}

		result[shortUrlId] = fullUrl
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
