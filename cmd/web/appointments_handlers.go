package main

import (
	"GoClinic/pkg/web/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Create Appointment function
func (app *application) createAppointment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		DateTime  string `json:"date_time"`
		DoctorID  int    `json:"doctor_id"`
		PatientID int    `json:"patient_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	appointment := &model.Appointment{
		DateTime:  input.DateTime,
		DoctorID:  input.DoctorID,
		PatientID: input.PatientID,
	}

	err = app.models.Appointments.Insert(appointment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"appointment": appointment}, nil)
}

// Get Appointments of the specific surgeon | function
func (app *application) getAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["appointmentId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such appointments, try another Appointment id")
		return
	}

	appointment, err := app.models.Appointments.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, appointment)
}

func (app *application) getSortedAppointments(w http.ResponseWriter, r *http.Request) {

	sortParam := r.URL.Query().Get("sort")
	sortDirection := r.URL.Query().Get("sort_direction")

	if sortDirection != "DESC" {
		sortDirection = "ASC"
	}

	filters := model.Filters{
		Sort:          sortParam,
		SortDirection: sortDirection,
		SortSafelist:  []string{"date_time", "doctor_id", "id", "-date_time", "-doctor_id", "-id"}, // Add any safe sorting criteria
	}

	// Call GetAllSortedByName method from the AppointmentModel instance
	appointments, err := app.models.Appointments.GetAllSortedByName(filters)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Marshal Appointments into JSON format
	jsonAppointments, err := json.Marshal(appointments)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	w.Write(jsonAppointments)
}

func (app *application) getFilteredAppointments(w http.ResponseWriter, r *http.Request) {

	filterParam := r.URL.Query().Get("filter")

	filtered_appointments, err := app.models.Appointments.GetFilteredByText(filterParam)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, filtered_appointments)
}

func (app *application) getPaginatedAppointments(w http.ResponseWriter, r *http.Request) {
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

	appointments, err := app.models.Appointments.GetPaginatedAppointments(limit, offset)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, appointments)
}

func (app *application) getAppointmentsOfDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["doctorId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such appointments, try another Doctor id")
		return
	}

	appointments, err := app.models.Appointments.Get_By_Doctor(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, appointments)
}

func (app *application) getAppointmentsOfPatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["patientId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such appointments, try another Patient id")
		return
	}

	appointments, err := app.models.Appointments.Get_By_Patient(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, appointments)
}

func (app *application) updateAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["appointmentId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Appointment Id written, try another")
		return
	}

	appointment, err := app.models.Appointments.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		DateTime  *string `json:"date_time"`
		DoctorID  *int    `json:"doctor_id"`
		PatientID *int    `json:"patient_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.DateTime != nil {
		appointment.DateTime = *input.DateTime
	}

	if input.DoctorID != nil {
		appointment.DoctorID = *input.DoctorID
	}

	if input.PatientID != nil {
		appointment.PatientID = *input.PatientID
	}

	err = app.models.Appointments.Update(appointment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error3")
		return
	}

	app.respondWithJSON(w, http.StatusOK, appointment)
}

func (app *application) deleteAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["appointmentId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "There is no such appointments with that Id")
		return
	}

	err = app.models.Appointments.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error4")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
