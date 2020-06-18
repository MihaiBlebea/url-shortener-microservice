package mysql

import (
	"context"

	"github.com/MihaiBlebea/url-shortener/shortener"

	// Mysql driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

// Find returns a redirect by code
func (r *Repository) Find(code string) (redirect shortener.Redirect, err error) {
	redirects, err := r.findBy("SELECT * FROM redirects WHERE code = ?", code)
	if err != nil {
		return redirects[0], err
	}
	return redirects[0], nil
}

// FindAll returns all redirects from the persistence layer
func (r *Repository) FindAll() (redirects []shortener.Redirect, err error) {
	return r.findBy("SELECT * FROM redirects WHERE id > 1")
}

func (r *Repository) findBy(sql string, params ...interface{}) (redirects []shortener.Redirect, err error) {
	ctx := context.Background()

	conn, err := r.client.Conn(ctx)
	if err != nil {
		return redirects, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return redirects, err
	}

	for rows.Next() {
		var redirect shortener.Redirect
		err = rows.Scan(
			&redirect.ID,
			&redirect.URL,
			&redirect.Code,
			&redirect.Hits,
			&redirect.Created,
			&redirect.Updated,
		)
		if err != nil {
			return redirects, err
		}
		redirects = append(redirects, redirect)
	}

	return redirects, nil
}
