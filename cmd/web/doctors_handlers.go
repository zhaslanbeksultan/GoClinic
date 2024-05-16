package main

import (
	"GoClinic/pkg/web/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Create Doctor function
func (app *application) createDoctor(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Speciality string `json:"speciality"`
		Phone      string `json:"phone"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	doctor := &model.Doctor{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Speciality: input.Speciality,
		Phone:      input.Phone,
	}

	err = app.models.Doctors.Insert(doctor)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"doctor": doctor}, nil)
}

// Get Doctors of the specific surgeon | function
func (app *application) getDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["doctorId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such doctors, try another Doctor id")
		return
	}

	doctor, err := app.models.Doctors.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, doctor)
}

func (app *application) getSortedDoctors(w http.ResponseWriter, r *http.Request) {

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

	// Call GetAllSortedByName method from the DoctorModel instance
	doctors, err := app.models.Doctors.GetAllSortedByName(filters)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Marshal doctors into JSON format
	jsonDoctors, err := json.Marshal(doctors)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	w.Write(jsonDoctors)
}

func (app *application) getFilteredDoctors(w http.ResponseWriter, r *http.Request) {

	filterParam := r.URL.Query().Get("filter")

	filtered_doctors, err := app.models.Doctors.GetFilteredByText(filterParam)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, filtered_doctors)
}

func (app *application) getPaginatedDoctors(w http.ResponseWriter, r *http.Request) {
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

	doctors, err := app.models.Doctors.GetPaginatedDoctors(limit, offset)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, doctors)
}

func (app *application) updateDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["doctorId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid doctor Id written, try another")
		return
	}

	doctor, err := app.models.Doctors.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		FirstName  *string `json:"first_name"`
		LastName   *string `json:"last_name"`
		Speciality *string `json:"speciality"`
		Phone      *string `json:"phone"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.FirstName != nil {
		doctor.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		doctor.LastName = *input.LastName
	}

	if input.Speciality != nil {
		doctor.Speciality = *input.Speciality
	}

	if input.Phone != nil {
		doctor.Phone = *input.Phone
	}

	err = app.models.Doctors.Update(doctor)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error3")
		return
	}

	app.respondWithJSON(w, http.StatusOK, doctor)
}

func (app *application) deleteDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["doctorId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such doctors with that Id")
		return
	}

	err = app.models.Doctors.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error4")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
