package data

import "time"

// Movie struct capital letter start means exported,
// to be visible in encoding/json package
// - used to hide form users, omitempty to hide if null, to leave json name ",omitempty"
// additionally, we can add ",string" to directive to convert int to string
type Movie struct {
	ID        int64     `json:"id"`                // Unique Integer ID for the movie
	CreatedAt time.Time `json:"-"`                 // Timestamps for when movies added to the database
	Title     string    `json:"title"`             // Movie title
	Year      int32     `json:"year,omitempty"`    // Movie release year
	Runtime   int32     `json:"runtime,omitempty"` // Movie Runtime (in minutes)
	Genres    []string  `json:"genres,omitempty"`  // Slice of genres for the movie (romance, comedy, etc.)
	Version   int32     `json:"version"`           // The version starts at number 1 and will be
	// incremented each time the movie information is updated
}
