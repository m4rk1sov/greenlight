package main

import (
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
