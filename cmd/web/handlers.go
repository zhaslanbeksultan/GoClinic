package main

import (
	"GoClinic/pkg/web/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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

// Create Registration function
func (app *application) createRegistration(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	registration := &model.Patient{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
	}

	err = app.models.Patients.Insert(registration)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"registration": registration}, nil)
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

//////////////////////////////////////////////////////////////////////////////////////////
