package main

import (
	"fmt"
	"net/http"
	"recipe_parrot/m/internal/models"
)

func (app *application) handler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /sms", app.handleIncomingMessage)

	return router
}

func (app *application) handleIncomingMessage(w http.ResponseWriter, r *http.Request) {
	defer w.Write([]byte(""))
	qs := r.URL.Query()
	phoneNumber := qs.Get("From")
	msgBody := qs.Get("Body")

	user, _ := app.models.Users.Get(phoneNumber)

	switch {
	case user != nil:
		app.handleAuthenticatedMessage(user, msgBody)
	case user == nil:
		app.handleUnauthenticatedMessage(phoneNumber, msgBody)
	}

	w.Header().Add("Content-Type", "text/xml")
	w.Write([]byte(""))
}

func (app *application) handleAuthenticatedMessage(user *models.User, msg string) {
	fmt.Printf("NUMBER: %s\nMESSAGE:%s\n", user.PhoneNumber, msg)
	switch msg {
	case "!setname":
		app.models.Users.AssignName(user.ID, "Noah Whitney")
	}
}

func (app *application) handleUnauthenticatedMessage(phoneNumber string, msg string) {
	switch msg {
	case "!register":
		fmt.Println("Registering User")
		user, err := app.models.Users.Create(phoneNumber)
		if err != nil {
			return
		}
		fmt.Printf("New user for %s\n", user.PhoneNumber)
	default:
		fmt.Printf("Unknown command: %s\n", msg)
	}
}
