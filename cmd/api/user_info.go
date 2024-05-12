package main

import (
	"errors"
	"greenlight.m4rk1sov.github.com/internal/data"
	"greenlight.m4rk1sov.github.com/internal/validator"
	"net/http"
)

// Vulnerable to user enumeration (meaning that attacker can know whether is user registered or not)
func (app *application) createUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name         string `json:"name"`
		Surname      string `json:"surname"`
		Email        string `json:"email"`
		PasswordHash string `json:"passwordHash"`
	}
	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.UserInfo{
		Name:      input.Name,
		Surname:   input.Surname,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.PasswordHash.Set(input.PasswordHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser2(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Insert the user data into the database.
	err = app.models.UserInfo.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Add the "user" permission for the new user.
	err = app.models.Permissions.ChangeRoleForUser(user.ID, "user")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//// After the user record has been created in the database, generate a new activation
	//// token for the user.
	//token, err := app.models.Tokens.New(user.ID, 40*time.Second, data.ScopeActivation)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}
	//app.background(func() {
	//	// As there are now multiple pieces of data that we want to pass to our email
	//	// templates, we create a map to act as a 'holding structure' for the data. This
	//	// contains the plaintext version of the activation token for the user, along
	//	// with their ID.
	//	data := map[string]any{
	//		"activationToken": token.Plaintext,
	//		"userID":          user.ID,
	//	}
	//	// Send the welcome email, passing in the map above as dynamic data.
	//	err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
	//	if err != nil {
	//		app.logger.PrintError(err, nil)
	//	}
	//})

	// Start the background email sending task
	app.background(func() {
		app.scheduleEmailSending(user)
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the plaintext token provided by the client.
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Retrieve the details of the user associated with the token using the
	// GetForToken() method (which we will create in a minute). If no matching record
	// is found, then we let the client know that the token they provided is not valid.
	user, err := app.models.UserInfo.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Update the user's activation status.
	user.Activated = true
	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.models.UserInfo.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}