package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Menu Singleton
	// Create a new menu
	v1.HandleFunc("/doctors", app.createDoctorHandler).Methods("POST")
	// Get a specific menu
	v1.HandleFunc("/doctors", app.getDoctorsHandler).Methods("GET")
	v1.HandleFunc("/doctors/{doctorId:[0-9]+}", app.requireActivatedUser(app.getDoctorHandler)).Methods("GET")
	// Update a specific menu
	v1.HandleFunc("/doctors/{doctorId:[0-9]+}", app.requirePermission("doctor.update", app.updateDoctorHandler)).Methods("PUT")
	// Delete a specific menu
	v1.HandleFunc("/doctors/{doctorId:[0-9]+}", app.requirePermission("doctor.delete", app.deleteDoctorHandler)).Methods("DELETE")

	v1.HandleFunc("/patients", app.createPatientHandler).Methods("POST")
	v1.HandleFunc("/patients", app.getPatientsHandler).Methods("GET")
	v1.HandleFunc("/patients/{patientId:[0-9]+}", app.requireActivatedUser(app.getPatientHandler)).Methods("GET")
	v1.HandleFunc("/patients/{patientId:[0-9]+}", app.updatePatientHandler).Methods("PUT")
	v1.HandleFunc("/patients/{patientId:[0-9]+}", app.deletePatientHandler).Methods("DELETE")

	v1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	v1.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

	log.Printf("Starting server on %d\n", app.config.port)
	//err := http.ListenAndServe(app.config.port, r)

	return app.authenticate(r)
}
