package main

import (
	"fmt"
	"greenlight.m4rk1sov.github.com/internal/data"
	"greenlight.m4rk1sov.github.com/internal/validator"
	"net/http"
	"time"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// anonymous struct subset of Movie struct will be in HTTP request body
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// initialize a json.Decode() instance to read from request
	// then Decode() to input in struct
	//err := json.NewDecoder(r.Body).Decode(&input)

	// Use the new readJSON() helper to decode the request body into the input struct.
	// If this returns an error we send the client the error message along with a 400
	// Bad Request status code, just like before.
	err := app.readJSON(w, r, &input)
	if err != nil {
		// use custom helper
		app.badRequestResponse(w, r, err)
		//app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Copy the values from the input struct to a new Movie struct.
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// content input struct
	fmt.Fprintf(w, "%+v\n", input)

	//fmt.Fprintln(w, "create")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// readIDParam in helpers.go
	id, err := app.readIDParam(r)
	if err != nil {
		// new helper
		app.notFoundResponse(w, r)
		//http.NotFound(w, r)
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
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// new helper
		app.serverErrorResponse(w, r, err)

		//app.logger.Print(err)
		//http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	// otherwise, interpolate Id in a placeholder response
	//fmt.Fprintf(w, "show the details of movie %d\n", id)
}
