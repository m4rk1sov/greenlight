package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// readIDParam in helpers.go
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// otherwise, interpolate Id in a placeholder response
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}
