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
	defer w.Write([]byte(""))
	qs := r.URL.Query()
	phoneNumber := qs.Get("From")
	msgBody := qs.Get("Body")

	fmt.Printf("message: %s\n", msgBody)

	user, err := app.models.Users.Get(phoneNumber)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.handleUnauthenticatedMessage(phoneNumber, msgBody)
		default:
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	}

	println("TEST2")
	res := app.handleAuthenticatedMessage(user, msgBody)

	fmt.Printf("RES: %s", res)
	w.Header().Add("Content-Type", "text/xml")
	_, err = w.Write([]byte(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}
func (app *application) handleAuthenticatedMessage(user *models.User, msg string) string {
	fmt.Printf("NUMBER: %s\nMESSAGE:%s\n", user.PhoneNumber, msg)
	msgStr := strings.Split(msg, " ")

	// Authenticated Commands
	if msgStr[0][0] == '!' {
		println("TESTR")
		cmd := strings.TrimPrefix(msgStr[0], "!")
		args := msgStr[1:]

		switch cmd {
		case "setname":
			err := app.models.Users.AssignName(user.ID, args[0], args[1])
			if err != nil {
				println(err.Error())
				res := twiml.MessagingMessage{
					Body: "Issue setting name, please try again",
					To:   user.PhoneNumber,
					From: app.config.twilio.fromNumber,
				}

				result, err := twiml.Messages([]twiml.Element{res})
				fmt.Printf("RESULT: %s\n", result)
				if err != nil {
					fmt.Println(err.Error())
				}
				return result
			}

			res := twiml.MessagingMessage{
				Body: fmt.Sprintf("Name successfully set to %s %s\n", args[0], args[1]),
				To:   user.PhoneNumber,
				From: app.config.twilio.fromNumber,
			}
			result, err := twiml.Messages([]twiml.Element{res})
			fmt.Printf("RESULT: %s\n", result)
			if err != nil {
				fmt.Println(err.Error())
			}
			return result
		case "echo":
			res := twiml.MessagingMessage{
				Body: "echo",
				To:   user.PhoneNumber,
				From: app.config.twilio.fromNumber,
			}
			result, err := twiml.Messages([]twiml.Element{res})
			fmt.Printf("RESULT: %s\n", result)
			if err != nil {
				fmt.Println(err.Error())
			}
			return result
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
	return ""
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
