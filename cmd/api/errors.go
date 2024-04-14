package main

import (
	"fmt"
	"net/http"
)

// generic helper logger for errors
func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

// the errorResponse() method to send JSON-format error with any type on message for versatility
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	// response with helper, if error occurs, returns 500 code status
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// the serverErrorResponse() method for unexpected problems at runtime 500 code
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// the notFoundResponse() method for 404 code
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// the methodNotAllowedResponse() method for 405
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this response", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
