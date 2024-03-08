package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/zhaslanbeksultan/GoClinic/pkg/web/model"
	"net/http"
	"strconv"
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
func (app *application) getAllRegistrations(w http.ResponseWriter, r *http.Request) {
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
