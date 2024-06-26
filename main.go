package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"recipe_parrot/m/internal/models"
	"time"

	_ "github.com/lib/pq"
)

// ARCHITECTURE
// WEB SCRAPER
// MESSAGES MODULE
// USER MODULE (PREFERENCES, ETC.)
// EVENT HANDLER
// LIST CREATOR
// DB MODEL
// CLI LAYER
// EVENT HANDLER is the main loop of program and awaits events/messages from the messages module.
// EVENT can be a recipe link to scrape, a user preference change, or a request for a list or other resource.
// All communication with user through MESSAGES MODULE.
// WEB SCRAPER handles scraping recipes from preconfigured sites.
// A recipe link sent through MESSAGES MODULE will hit EVENT HANDLER and be sent to WEB SCRAPER which will send a recipe back
// to EVENT HANDLER to then be stored in DB MODEL for later use.

type config struct {
	port   int
	twilio struct {
		accountSid string
		authToken  string
		fromNumber string
	}
	db struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	models models.Models
}

func main() {
	var cfg config

	// server
	flag.IntVar(&cfg.port, "server port", 6000, "port for recipe parrot server")

	// Database Config
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://recipe_parrot_user:Noah2002ndw@localhost/recipe_parrot?sslmode=disable",
		"DB connection string")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m",
		"PostgreSQL max connection idle time")

	// twilio
	flag.StringVar(&cfg.twilio.accountSid, "twilio account sid", os.Getenv("TWILIO_ACCOUNT_SID"), "account sid for twilio messaging api")
	flag.StringVar(&cfg.twilio.authToken, "twilio auth token", os.Getenv("TWILIO_AUTH_TOKEN"), "auth token for twilio messaging api")
	flag.StringVar(&cfg.twilio.fromNumber, "twilio from number", "+12402402001", "phone number used to send and receive messages to server")

	db, err := openDB(cfg)
	if err != nil {
		return
	}

	app := &application{
		config: cfg,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.handler(),
	}

	fmt.Printf("Server started on port %d\n", app.config.port)
	err = srv.ListenAndServe()
	fmt.Print(err.Error())

	// client := twilio.NewRestClient()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
