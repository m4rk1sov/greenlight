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
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	// Add the route for the PUT /v1/movies/:id endpoint.
	// Require a PATCH request, rather than PUT.
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	// Add the route for the DELETE /v1/movies/:id endpoint.
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))

	router.HandlerFunc(http.MethodGet, "/v1/modules", app.requireRole("user", app.listModulesInfoHandler))
	router.HandlerFunc(http.MethodPost, "/v1/modules", app.requireRole("admin", app.createModuleInfoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/modules/:id", app.requireRole("user", app.getModuleInfoHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/modules/:id", app.requireRole("admin", app.editModuleInfoHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/modules/:id", app.requireRole("admin", app.deleteModuleInfoHandler))

	router.HandlerFunc(http.MethodPost, "/v1/departments", app.requireRole("admin", app.createDepartmentInfoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/departments/:id", app.requireRole("user", app.getDepartmentInfoHandler))

	// Add the route for the POST /v1/userinfo endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/userinfo", app.createUserInfoHandler)
	// Add the route for the PUT /v1/userinfo/activated endpoint.
	router.HandlerFunc(http.MethodPut, "/v1/userinfo/activated", app.activateUserInfoHandler)
	// Add the route for the POST /v1/tokens/authentication endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/tokens2/authentication", app.createAuthenticationTokenHandler2)

	// Add the route for the POST /v1/users endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	// Add the route for the PUT /v1/users/activated endpoint.
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	// Add the route for the POST /v1/tokens/authentication endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Wrap the router with the rateLimit() middleware.
	// Use the authenticate() middleware on all requests.

	// change to authenticate(router) for book
	return app.recoverPanic(app.rateLimit(app.authenticate2(router)))
}
