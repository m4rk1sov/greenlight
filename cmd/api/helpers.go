package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// retrieve Id convert it to integer and return, otherwise return 0, error
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	// get the value of ID from slice, Id is int but byName() returns String,
	// so we need to convert it to base 10 integer,
	// bit size 64, otherwise Id is invalid
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// define envelope type
type envelope map[string]any

// define writeJSON helper for sending responses
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// encode the data to JSON
	// use the json.MarshalIndent to add whitespace and tab indent
	// Unfortunately it takes more memory and slower than a regular one
	// if application relies on high amount of traffic need to switch to original
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	// we know that there are no more errors before writing to the response
	// safe to add the headers in the map, even nil
	for key, value := range headers {
		w.Header()[key] = value
	}

	//add the Content-Type: application/json header then status code and JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
