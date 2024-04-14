package main

import (
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// create a map which holds the response info
	//data := map[string]string{
	//	"status":      "available",
	//	"environment": app.config.env,
	//	"version":     version,
	//}

	// declare envelop map for data, nested values
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		// use new helper
		app.serverErrorResponse(w, r, err)

		//app.logger.Print(err)
		//http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	//// append a new line to json, make it easier to read in the terminal
	//js = append(js, '\n')
	//
	//// fixed JSON response, raw string literal for double quotes without escaping them
	//// we also use %q to interpolate data into double quotes
	//// js := `{"status": "available", "environment": %q, "version": %q}`
	//// js = fmt.Sprintf(js, app.config.env, version)
	//
	//// setting Content-type to application/json instead of the default one
	//// encoding works without problem, safely set necessary HTTP headers
	//w.Header().Set("Content-Type", "application/json")
	//
	//// write the json []byte
	//w.Write(js)
	//
	//// hardcoded content
	////fmt.Fprintln(w, "status: available")
	////fmt.Fprintf(w, "environment: %s\n", app.config.env)
	////fmt.Fprintf(w, "version: %s\n", version)
}
