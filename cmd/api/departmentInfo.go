package main

import (
	"errors"
	"fmt"
	"greenlight.m4rk1sov.github.com/internal/data"
	"greenlight.m4rk1sov.github.com/internal/validator"
	"net/http"
)

func (app *application) createDepartmentInfoHandler(w http.ResponseWriter, r *http.Request) {
	// anonymous struct subset of Movie struct will be in HTTP request body
	var input struct {
		DepartmentName     string `json:"departmentName"`
		StaffQuantity      int64  `json:"staffQuantity"`
		DepartmentDirector string `json:"departmentDirector"`
		Module_Info        int64  `json:"module_Info"`
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
	// Note that the movie variable contains a *pointer* to a Movie struct.
	departmentInfo := &data.DepartmentInfo{
		DepartmentName:     input.DepartmentName,
		StaffQuantity:      input.StaffQuantity,
		DepartmentDirector: input.DepartmentDirector,
		Module_Info:        input.Module_Info,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateDepartment(v, departmentInfo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our movies model, passing in a pointer to the
	// validated movie struct. This will create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.DepartmentInfo.Insert(departmentInfo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//// content input struct
	//fmt.Fprintf(w, "%+v\n", input)

	//fmt.Fprintln(w, "create")

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/departments/%d", departmentInfo.ID))
	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"departments": departmentInfo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getDepartmentInfoHandler(w http.ResponseWriter, r *http.Request) {
	// readIDParam in helpers.go
	id, err := app.readIDParam(r)
	if err != nil {
		// new helper
		app.notFoundResponse(w, r)
		//http.NotFound(w, r)
		return
	}

	//// Create a new instance of Movie struct, containing the ID we extracted
	//// from the URL and some dummy data. Also notice that we deliberately
	//// haven's set a value for the Year field
	//movie := data.Movie{
	//	ID:        id,
	//	CreatedAt: time.Now(),
	//	Title:     "Casablanca",
	//	Runtime:   102,
	//	Genres:    []string{"drama", "romance", "war"},
	//	Version:   1,
	//}

	// Call the Get() method to fetch the data for a specific movie. We also need to
	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.

	departmentInfo, err := app.models.DepartmentInfo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct JSON and send it as the HTTP response
	err = app.writeJSON(w, http.StatusOK, envelope{"departmentInfo": departmentInfo}, nil)
	if err != nil {
		// new helper
		app.serverErrorResponse(w, r, err)

		//app.logger.Print(err)
		//http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	// otherwise, interpolate Id in a placeholder response
	//fmt.Fprintf(w, "show the details of movie %d\n", id)
}
