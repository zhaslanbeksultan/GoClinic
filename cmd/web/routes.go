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

	// Create a new menu
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

	users1 := r.PathPrefix("/api/v1").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
