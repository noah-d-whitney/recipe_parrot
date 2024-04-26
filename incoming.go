package main

import (
	"errors"
	"fmt"
	"net/http"
	"recipe_parrot/m/internal/models"
	"strings"

	"github.com/twilio/twilio-go/twiml"
)

func (app *application) handler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /sms", app.handleIncomingMessage)

	return router
}

func (app *application) handleIncomingMessage(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	phoneNumber := qs.Get("From")
	msgBody := qs.Get("Body")

	fmt.Printf("message: %s\n", msgBody)

	user, err := app.models.Users.Get(phoneNumber)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			fmt.Printf("UNAUTH")
			app.handleUnauthenticatedMessage(phoneNumber, msgBody)
		default:
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	}

	println("TEST2")
	res := app.handleAuthenticatedMessage(user, msgBody)

	message := twiml.MessagingMessage{
		Body: res,
		To:   phoneNumber,
		From: app.config.twilio.fromNumber,
	}

	result, err := twiml.Messages([]twiml.Element{message})
	if err != nil {
		return
	}

	fmt.Printf("RES: %s", res)
	w.Header().Add("Content-Type", "text/xml")
	_, err = w.Write([]byte(result))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (app *application) handleAuthenticatedMessage(user *models.User, msg string) string {
	msgStr := strings.Split(msg, " ")

	switch msgStr[0] {
	case "LIST":
		return "Send list here!!!"
	case "NEW":
		return "Create new shopping trip"
	}
	return ""
}

func (app *application) handleUnauthenticatedMessage(phoneNumber string, msg string) {
	switch msg {
	case "REGISTER":
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
