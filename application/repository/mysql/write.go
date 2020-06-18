package mysql

import (
	"context"
	"errors"

	"github.com/MihaiBlebea/url-shortener/shortener"

	// Mysql driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

// Migrate generates a users table
func (r *Repository) Migrate() error {
	ctx := context.Background()

	conn, err := r.client.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	schema := `CREATE TABLE IF NOT EXISTS redirects(
		id int NOT NULL AUTO_INCREMENT,
		url VARCHAR(250),
		code VARCHAR(100) UNIQUE,
		hits VARCHAR(100),
		created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id));`
	stmt, err := conn.PrepareContext(ctx, schema)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

// Store adds a redirect model to the db
func (r *Repository) Store(redirect shortener.Redirect) error {
	ctx := context.Background()

	conn, err := r.client.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(
		ctx,
		`INSERT INTO redirects 
			(url, code, hits, created) 
		VALUES 
			(?, ?, ?, ?)`,
	)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		redirect.URL,
		redirect.Code,
		redirect.Hits,
		redirect.Created,
	)
	return err
}

// Update updates a redirect
func (r *Repository) Update(redirect shortener.Redirect) error {
	if redirect.ID == 0 {
		return errors.New("Could not update user as it doesn't have an id")
	}

	ctx := context.Background()

	conn, err := r.client.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(
		ctx,
		`UPDATE 
			redirects 
		SET 
			url = ?, code = ?, hits = ? 
		WHERE 
			id=?`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		redirect.URL,
		redirect.Code,
		redirect.Hits,
		redirect.ID,
	)
	return err
}
