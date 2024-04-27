package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrSiteNotSupported = errors.New("site is unsupported for scraping")
)

type Site string

const (
	all_recipes Site = "allrecipes"
)

type SiteModel struct {
	db *sql.DB
}

func (m *SiteModel) parse(url string) (Site, error) {
	s := url
	s = strings.TrimPrefix(s, "https://www.")
	idx := strings.Index(s, ".")
	s = s[:idx]
	fmt.Printf("SITE STRING: %s\n", s)

	stmt := `SELECT name FROM sites WHERE name = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result Site
	err := m.db.QueryRowContext(ctx, stmt, s).Scan(&result)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Site(""), ErrSiteNotSupported
		default:
			return Site(""), err
		}
	}
	fmt.Printf("RESULT SITE: %s\n", result)

	return result, nil
}

func (m *SiteModel) Scrape(url string) (*Recipe, error) {
	site, err := m.parse(url)
	if err != nil {
		return nil, err
	}

	switch site {
	case all_recipes:
		r, err := scrape_all_recipes(url)
		if err != nil {
			return nil, err
		}
		return r, nil
	default:
		return nil, ErrSiteNotSupported
	}
}
