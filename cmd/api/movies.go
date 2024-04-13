package main

import (
	"fmt"
	"greenlight.m4rk1sov.github.com/internal/data"
	"net/http"
	"time"
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

	// Create a new instance of Movie struct, containing the ID we extracted
	// from the URL and some dummy data. Also notice that we deliberately
	// haven's set a value for the Year field
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// Encode the struct JSON and send it as the HTTP response
	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
	// otherwise, interpolate Id in a placeholder response
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}
