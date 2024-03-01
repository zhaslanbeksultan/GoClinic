package main

import (
	"encoding/json"
	"github.com/zhaslanbeksultan/GoClinic/pkg/web/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
//Create Registration function
func (app *application) createRegistration(w http.ResponseWriter, r *http.Request) {
	var input struct {
		First_name          string `json:"first_name"`
		Last_name    string `json:"last_name"`
		Phone uint   `json:"phone"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	registration := &model.Patient{
		First_name:         input.First_name,
		Last_name:    		input.Last_name,
		Phone: 				input.Phone,
	}

	err = app.models.Registrations.Insert(registration)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, registration)
}

//Get Registrations of the secific surgeon | function
func (app *application) getAllRegistrations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["registrationId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such doctor in our clinic, try another id")
		return
	}

	registration, err := app.models.Registrations.Get(id)
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
		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}

	menu, err := app.models.Menus.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Title          *string `json:"title"`
		Description    *string `json:"description"`
		NutritionValue *uint   `json:"nutritionValue"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		menu.Title = *input.Title
	}

	if input.Description != nil {
		menu.Description = *input.Description
	}

	if input.NutritionValue != nil {
		menu.NutritionValue = *input.NutritionValue
	}

	err = app.models.Menus.Update(menu)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, menu)
}