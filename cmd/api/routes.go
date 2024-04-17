package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	// router instance
	router := httprouter.New()

	// convert our own helpers to http.Handler 404 code error using adapter SDP beyba xD
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// likewise, convert to 405 error, basically making custom which is supported by http.Handler
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// relevant methods
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	// Add the route for the PUT /v1/movies/:id endpoint.
	// Require a PATCH request, rather than PUT.
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	// Add the route for the DELETE /v1/movies/:id endpoint.
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	// return instance
	return router
}
