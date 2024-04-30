package main

import (
	"encoding/json"
	"errors"
	"fmt"
	pkg "github.com/zhaslanbeksultan/GoClinic/pkg/web/model"
	"github.com/zhaslanbeksultan/GoClinic/pkg/web/validator"
	"net/http"
	"time"

	"strconv"

	"github.com/gorilla/mux"
)

// patient handlers
func (app *application) getPatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["patientId"])
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid patient number")
		return
	}

	patient, err := app.models.Patients.GetPatient(id)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Patient not found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, patient)
}
func (app *application) getPatientsHandler(w http.ResponseWriter, r *http.Request) {
	patients, err := app.models.Patients.GetPatients()

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, patients)

}
func (app *application) createPatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient pkg.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if patient.FirstName == "" || patient.LastName == "" {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = app.models.Patients.InsertPatient(&patient)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, patient)
}

func (app *application) updatePatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient pkg.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if patient.FirstName == "" || patient.LastName == "" {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = app.models.Patients.UpdatePatient(&patient)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, patient)
}

func (app *application) deletePatientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["patientId"])

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	err = app.models.Patients.DeletePatient(id)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusNoContent, nil)
}

// Doctor Handlers
func (app *application) getDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName  string
		SecondName string
		pkg.Filters
	}
	// Initialize a new Validator instance.
	v := validator.New()
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.FirstName = app.readString(qs, "firstName", "")
	input.SecondName = app.readString(qs, "secondName", "")
	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "doctorId")
	input.Filters.SortSafelist = []string{"doctorId", "firstName", "secondName", "speciality", "phone", "-doctorId", "-firstName", "-secondName", "-speciality", "-phone"}
	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send the client a response if necessary.
	if pkg.ValidateFilters(v, input.Filters); !v.Valid() {
		//app.failedValidationResponse(w, r, v.Errors)
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	if !v.Valid() {
		//app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Dump the contents of the input struct in a HTTP response.
	//fmt.Fprintf(w, "%+v\n", input)

	doctors, metadata, err := app.models.Doctors.GetDoctors(input.FirstName, input.SecondName, input.Filters)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	fmt.Println(metadata)
	app.respondWithJSON(w, http.StatusOK, doctors)
}
func (app *application) getDoctorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["doctorId"])
	fmt.Println(id)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	doctor, err := app.models.Doctors.GetDoctor(id)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Doctor not found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, doctor)
}
func (app *application) createDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var doctor pkg.Doctor
	err := json.NewDecoder(r.Body).Decode(&doctor)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if doctor.FirstName == "" || doctor.SecondName == "" {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = app.models.Doctors.InsertDoctor(&doctor)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, doctor)
}

func (app *application) updateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var doctor pkg.Doctor
	err := json.NewDecoder(r.Body).Decode(&doctor)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if doctor.FirstName == "" || doctor.SecondName == "" {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = app.models.Doctors.UpdateDoctor(&doctor)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, doctor)
}

func (app *application) deleteDoctorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["doctorId"])

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	err = app.models.Doctors.DeleteDoctor(id)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusNoContent, nil)
}

// User Handlers
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
		return
	}
	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.

	user := &pkg.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if pkg.ValidateUser(v, user); !v.Valid() {
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
		case errors.Is(err, pkg.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			//app.failedValidationResponse(w, r, v.Errors)
		default:
			//app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.models.Permissions.AddForUser(user.ID, "doctor.read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, pkg.ScopeActivation)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
		return
	}
	var res struct {
		Token *string   `json:"token"`
		User  *pkg.User `json:"user"`
	}

	res.Token = &token.Plaintext
	res.User = user

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": res}, nil)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
	}

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		//app.badRequestResponse(w, r, err)
		return
	}
	// Validate the plaintext token provided by the client.
	v := validator.New()
	if pkg.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		//app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Retrieve the details of the user associated with the token using the
	// GetForToken() method (which we will create in a minute). If no matching record
	// is found, then we let the client know that the token they provided is not valid.
	user, err := app.models.Users.GetForToken(pkg.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, err):
			v.AddError("token", "invalid or expired activation token")
			//app.failedValidationResponse(w, r, v.Errors)
		default:
			//app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Update the user's activation status.
	user.Activated = true
	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, err):
			//app.editConflictResponse(w, r)
		default:
			//app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.models.Tokens.DeleteAllForUser(pkg.ScopeActivation, user.ID)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
	}
}
