package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/zhaslanbeksultan/GoClinic/pkg/web/model"
	"github.com/zhaslanbeksultan/GoClinic/validator"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type envelope map[string]interface{}

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error1")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists this
	// will return the empty string "".
	s := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}
	// Otherwise return the string.
	return s
}

// The readCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if csv == "" {
		return defaultValue
	}
	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Extract the value from the query string.
	s := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}
	// Try to convert the value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	// Otherwise, return the converted integer value.
	return i
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope,
	headers http.Header) error {
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded JSON. Use
	// no line prefix and tab indents for each element.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// At this point, we know that we won't encounter any more errors before writing the response,
	// so it's safe to add any headers that we want to include. We loop through the header map
	// and add each header to the http.ResponseWriter header map. Note that it's OK if the
	// provided header map is nil. Go doesn't through an error if you try to range over (
	// or generally, read from) a nil map
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(js); err != nil {
		//app.logger.PrintError(err, nil)
		return err
	}

	return nil
}

// Create Registration function
func (app *application) createRegistration(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	registration := &model.Patient{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
	}

	err = app.models.Patients.Insert(registration)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error2")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, registration)
}

// Get Registrations of the specific surgeon | function
func (app *application) getRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["registrationId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such registrations, try another Registration id")
		return
	}

	registration, err := app.models.Patients.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, registration)
}

func (app *application) getSortedRegistrations(w http.ResponseWriter, r *http.Request) {

	sortParam := r.URL.Query().Get("sort")
	sortDirection := r.URL.Query().Get("sort_direction")

	if sortDirection != "DESC" {
		sortDirection = "ASC"
	}

	filters := model.Filters{
		Sort:          sortParam,
		SortDirection: sortDirection,
		SortSafelist:  []string{"first_name", "last_name", "id", "-first_name", "-last_name", "-id"}, // Add any safe sorting criteria
	}

	// Call GetAllSortedByName method from the PatientModel instance
	registrations, err := app.models.Patients.GetAllSortedByName(filters)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Marshal registrations into JSON format
	jsonRegistrations, err := json.Marshal(registrations)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	w.Write(jsonRegistrations)
}

func (app *application) getFilteredRegistrations(w http.ResponseWriter, r *http.Request) {

	filterParam := r.URL.Query().Get("filter")

	filtered_registrations, err := app.models.Patients.GetFilteredByText(filterParam)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, filtered_registrations)
}

func (app *application) getPaginatedRegistrations(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid limit data")
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid offset data")
		return
	}

	patients, err := app.models.Patients.GetPaginatedPatients(limit, offset)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, patients)
}

func (app *application) updateRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["registrationId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid registration Id written, try another")
		return
	}

	registration, err := app.models.Patients.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Phone     *string `json:"phone"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.FirstName != nil {
		registration.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		registration.LastName = *input.LastName
	}

	if input.Phone != nil {
		registration.Phone = *input.Phone
	}

	err = app.models.Patients.Update(registration)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error3")
		return
	}

	app.respondWithJSON(w, http.StatusOK, registration)
}

func (app *application) deleteRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["registrationId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such registrations with that Id")
		return
	}

	err = app.models.Patients.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error4")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		//app.badRequestResponse(w, r, err)
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.

	user := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload: password")
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if model.ValidateUser(v, user); !v.Valid() {
		//app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, model.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			//app.failedValidationResponse(w, r, v.Errors)
		default:
			//app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.models.Permissions.AddForUser(user.ID, "club.read")
	if err != nil {
		//app.serverErrorResponse(w, r, err)
		return
	}

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
		return
	}
	var res struct {
		Token *string     `json:"token"`
		User  *model.User `json:"user"`
	}

	res.Token = &token.Plaintext
	res.User = user

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": res}, nil)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
	}

}
