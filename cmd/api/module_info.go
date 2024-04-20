package main

import (
	"errors"
	"fmt"
	"greenlight.m4rk1sov.github.com/internal/data"
	"greenlight.m4rk1sov.github.com/internal/validator"
	"net/http"
)

func (app *application) createModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	// anonymous struct subset of Movie struct will be in HTTP request body
	var input struct {
		ModuleName     string       `json:"moduleName"`
		ModuleDuration data.Runtime `json:"moduleDuration"`
		ExamType       string       `json:"examType"`
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
	module_info := &data.Module_info{
		ModuleName:     input.ModuleName,
		ModuleDuration: input.ModuleDuration,
		ExamType:       input.ExamType,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateModule(v, module_info); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our movies model, passing in a pointer to the
	// validated movie struct. This will create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.Module_info.Insert(module_info)
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
	headers.Set("Location", fmt.Sprintf("/v1/modules/%d", module_info.ID))
	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"modules": module_info}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	module_info, err := app.models.Module_info.Get(id)
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
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": module_info}, nil)
	if err != nil {
		// new helper
		app.serverErrorResponse(w, r, err)

		//app.logger.Print(err)
		//http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	// otherwise, interpolate Id in a placeholder response
	//fmt.Fprintf(w, "show the details of movie %d\n", id)
}

func (app *application) editModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	module_info, err := app.models.Module_info.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.
	// Use pointers for the Title, Year and Runtime fields.
	var input struct {
		ModuleName     *string       `json:"moduleName"`
		ModuleDuration *data.Runtime `json:"moduleDuration,omitempty"`
		ExamType       *string       `json:"examType"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// If the input.Title value is nil then we know that no corresponding "title" key/
	// value pair was provided in the JSON request body. So we move on and leave the
	// movie record unchanged. Otherwise, we update the movie record with the new title
	// value. Importantly, because input.Title is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our movie record.
	if input.ModuleName != nil {
		module_info.ModuleName = *input.ModuleName
	}
	// We also do the same for the other fields in the input struct.
	if input.ModuleDuration != nil {
		module_info.ModuleDuration = *input.ModuleDuration
	}
	if input.ExamType != nil {
		module_info.ExamType = *input.ExamType
	}

	//// Copy the values from the request body to the appropriate fields of the movie
	//// record.
	//movie.Title = input.Title
	//movie.Year = input.Year
	//movie.Runtime = input.Runtime
	//movie.Genres = input.Genres

	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateModule(v, module_info); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//// Pass the updated movie record to our new Update() method.
	//err = app.models.Movies.Update(movie)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}

	// Intercept any ErrEditConflict error and call the new editConflictResponse()
	// helper.
	err = app.models.Module_info.Update(module_info)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": module_info}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the movie from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Module_info.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Module successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listModulesInfoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ModuleName string
		ExamType   string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.ModuleName = app.readString(qs, "moduleName", "")
	input.ExamType = app.readString(qs, "examType", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "moduleName", "moduleDuration", "-id", "-moduleName", "-moduleDuration"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the GetAll() method to retrieve the movies, passing in the various filter
	// parameters.
	// Accept the metadata struct as a return value.
	modules_info, metadata, err := app.models.Module_info.GetAllModules(input.ModuleName, input.ExamType, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	// Include the metadata in the response envelope.
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": modules_info, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
