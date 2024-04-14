package data

import (
	"greenlight.m4rk1sov.github.com/internal/validator"
	"time"
)

// Movie struct capital letter start means exported,
// to be visible in encoding/json package
// - used to hide form users, omitempty to hide if null, to leave json name ",omitempty"
// additionally, we can add ",string" to directive to convert int to string
type Movie struct {
	ID        int64     `json:"id"`             // Unique Integer ID for the movie
	CreatedAt time.Time `json:"-"`              // Timestamps for when movies added to the database
	Title     string    `json:"title"`          // Movie title
	Year      int32     `json:"year,omitempty"` // Movie release year
	// now Runtime type, if null, skips the method
	Runtime Runtime  `json:"runtime,omitempty"` // Movie Runtime (in minutes)
	Genres  []string `json:"genres,omitempty"`  // Slice of genres for the movie (romance, comedy, etc.)
	Version int32    `json:"version"`           // The version starts at number 1 and will be
	// incremented each time the movie information is updated
}

// To prevent duplication, we can collect the validation checks for a movie into a standalone
// ValidateMovie() function
func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
