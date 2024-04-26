package main

import (
	"fmt"
	"net/http"
)

func (app *application) invalidInputResponse(_ http.ResponseWriter, input string) {
	fmt.Printf("invalid input: %s\n", input)
}
