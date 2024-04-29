package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/zhaslanbeksultan/GoClinic/pkg/web/jsonlog"
	"github.com/zhaslanbeksultan/GoClinic/pkg/web/model"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {

	var cfg config
	flag.StringVar(&cfg.port, "port", ":8080", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:$F00tba11@localhost:5432/postgres?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Create a new menu
	v1.HandleFunc("/creation", app.createRegistration).Methods("POST")
	// Get a specific patient
	v1.HandleFunc("/registrations/{registrationId:[0-9]+}", app.getRegistration).Methods("GET")
	// Update a specific patient
	v1.HandleFunc("/registrations/{registrationId:[0-9]+}", app.updateRegistration).Methods("PUT")
	// // Delete a specific patient
	v1.HandleFunc("/registrations/{registrationId:[0-9]+}", app.deleteRegistration).Methods("DELETE")
	// Get sorted patients list
	v1.HandleFunc("/registrations/sorting", app.getSortedRegistrations).Methods("GET")
	// Get filtered patients list
	v1.HandleFunc("/registrations", app.getFilteredRegistrations).Methods("GET")
	// Get paginated patients list
	v1.HandleFunc("/registrations/paginated", app.getPaginatedRegistrations).Methods("GET")

	users1 := r.PathPrefix("/api/v1").Subrouter()
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open(`postgres`, cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
