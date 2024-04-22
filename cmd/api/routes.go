package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
	// router instance
	router := httprouter.New()

	// convert our own helpers to http.Handler 404 code error using adapter SDP beyba xD
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// likewise, convert to 405 error, basically making custom which is supported by http.Handler
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// router for healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// relevant methods
	// Add the route for the GET /v1/movies endpoint.
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	// Add the route for the PUT /v1/movies/:id endpoint.
	// Require a PATCH request, rather than PUT.
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	// Add the route for the DELETE /v1/movies/:id endpoint.
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)

	router.HandlerFunc(http.MethodGet, "/v1/modules", app.listModulesInfoHandler)
	router.HandlerFunc(http.MethodPost, "/v1/modules", app.createModuleInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/modules/:id", app.getModuleInfoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/modules/:id", app.editModuleInfoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/modules/:id", app.deleteModuleInfoHandler)

	router.HandlerFunc(http.MethodPost, "/v1/departments", app.createDepartmentInfoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/departments/:id", app.getDepartmentInfoHandler)

	// Add the route for the POST /v1/users endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	// Wrap the router with the rateLimit() middleware.
	return app.recoverPanic(app.rateLimit(router))
}
