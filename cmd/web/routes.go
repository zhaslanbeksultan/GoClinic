package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Create a new patient
	v1.HandleFunc("/creation", app.requireActivatedUser(app.createRegistration)).Methods("POST")
	// Get a specific patient
	v1.HandleFunc("/registrations/{registrationId:[0-9]+}", app.requireActivatedUser(app.getRegistration)).Methods("GET")
	// Update a specific patient
	v1.HandleFunc("/registrations/{registrationId:[0-9]+}", app.requireActivatedUser(app.updateRegistration)).Methods("PUT")
	// // Delete a specific patient
	v1.HandleFunc("/registrations/{registrationId:[0-9]+}", app.requirePermissions("patient.delete", app.deleteRegistration)).Methods("DELETE")
	// Get sorted patients list
	v1.HandleFunc("/registrations/sorting", app.requireActivatedUser(app.getSortedRegistrations)).Methods("GET")
	// Get filtered patients list
	v1.HandleFunc("/registrations", app.requireActivatedUser(app.getFilteredRegistrations)).Methods("GET")
	// Get paginated patients list
	v1.HandleFunc("/registrations/paginated", app.requireActivatedUser(app.getPaginatedRegistrations)).Methods("GET")
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	v2 := r.PathPrefix("/api/v1").Subrouter()

	// Create a new doctor
	v2.HandleFunc("/doctorcreation", app.requireActivatedUser(app.createDoctor)).Methods("POST")
	// Get a specific doctor
	v2.HandleFunc("/doctors/{doctorId:[0-9]+}", app.requireActivatedUser(app.getDoctor)).Methods("GET")
	// Update a specific doctor
	v2.HandleFunc("/doctors/{doctorId:[0-9]+}", app.requireActivatedUser(app.updateDoctor)).Methods("PUT")
	// // Delete a specific doctor
	v2.HandleFunc("/doctors/{doctorId:[0-9]+}", app.requirePermissions("doctor.delete", app.deleteDoctor)).Methods("DELETE")
	// Get sorted doctors list
	v2.HandleFunc("/doctors/sorting", app.requireActivatedUser(app.getSortedDoctors)).Methods("GET")
	// Get filtered doctors list
	v2.HandleFunc("/doctors", app.requireActivatedUser(app.getFilteredDoctors)).Methods("GET")
	// Get paginated doctors list
	v2.HandleFunc("/doctors/paginated", app.requireActivatedUser(app.getPaginatedDoctors)).Methods("GET")
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	v3 := r.PathPrefix("/api/v1").Subrouter()

	// Create a new appointment
	v3.HandleFunc("/newappointment", app.requireActivatedUser(app.createAppointment)).Methods("POST")
	// Get a specific appointment
	v3.HandleFunc("/appointments/{appointmentId:[0-9]+}", app.requireActivatedUser(app.getAppointment)).Methods("GET")
	// Get filtered appointment list
	v3.HandleFunc("/appointments", app.requireActivatedUser(app.getFilteredAppointments)).Methods("GET")
	// Update a specific appointment
	v3.HandleFunc("/appointments/{appointmentId:[0-9]+}", app.requireActivatedUser(app.updateAppointment)).Methods("PUT")
	// Delete a specific appointment
	v3.HandleFunc("/appointments/{appointmentsId:[0-9]+}", app.requirePermissions("appointment.delete", app.deleteAppointment)).Methods("DELETE")
	// Get sorted doctors list
	v2.HandleFunc("/appointments/sorting", app.requireActivatedUser(app.getSortedAppointments)).Methods("GET")
	// Get paginated appointments list
	v3.HandleFunc("/appointments/paginated", app.requireActivatedUser(app.getPaginatedAppointments)).Methods("GET")
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	users1 := r.PathPrefix("/api/v1").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
