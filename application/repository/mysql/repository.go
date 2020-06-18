package mysql

import (
	"database/sql"
	"fmt"

	// Mysql driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

// Repository struct
type Repository struct {
	client *sql.DB
}

// NewClient creates a new client to the mysql server
func NewClient(username, password, host, port, database string) (client *sql.DB, err error) {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username,
		password,
		host,
		port,
		database,
	)
	client, err = sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewRepository creates a new repository
func NewRepository(client *sql.DB) (repository *Repository, err error) {
	repository = &Repository{client}
	err = repository.Migrate()
	if err != nil {
		return repository, err
	}
	return repository, nil
}
