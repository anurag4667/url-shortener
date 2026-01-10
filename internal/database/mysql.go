package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQL(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) Save(id, url string) error {
	log.Println(">>> ENTERED MYSQL SAVE <<<")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO urls (id, original_url) VALUES (?, ?)",
		id, url,
	)
	return err
}

func (s *MySQLStore) Get(id string) (string, bool, error) {
	var url string
	err := s.db.QueryRow(
		"SELECT original_url FROM urls WHERE id = ?",
		id,
	).Scan(&url)

	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return url, true, nil
}

func (s *MySQLStore) IncrementClicks(id string) {
	// fire-and-forget analytics
	go s.db.Exec(
		"UPDATE urls SET click_count = click_count + 1 WHERE id = ?",
		id,
	)
}
