package main

import (
	"errors"
	"fmt"
	"net/http"
	"recipe_parrot/m/internal/models"

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
	switch msg {
	case "LIST":
		list, err := app.models.Lists.GetCurrentList(user.ID)
		if err != nil {
			fmt.Print(err.Error())
			return "Issue getting your list right now! Please try again soon"
		}
		ingredientsList := list.GetIngredientsList()
		response, err := ingredientsList.GenerateMessage()
		if err != nil {
			fmt.Print(err.Error())
			return "Something went wrong, please try again soon"
		}
		return response
	case "NEW":
		err := app.models.Lists.StartNewList(user.ID)
		if err != nil {
			fmt.Print(err.Error())
			return "Issue starting new list! Please try again soon"
		}
		return "New shopping list started. Send HELP for help."
	default:
		recipe, err := app.models.Sites.Scrape(msg)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrSiteNotSupported):
				return "Recipes from provided site are currently not supported"
			default:
				return "Unrecognized input, send HELP for help"
			}
		}
		recipe.UserID = user.ID

		err = app.models.Recipes.Create(recipe)
		if err != nil {
			return "Issue saving recipe, please try again later"
		}

		return "Recipe added to shopping trip, send LIST for current list"
	}
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
