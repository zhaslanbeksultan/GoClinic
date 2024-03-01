package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Patient struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	First_name string `json:"firstName"`
	Last_name  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type PatientModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m PatientModel) Insert(patient *Patient) error {
	// Insert a new menu item into the database.
	query := `
		INSERT INTO patients (first_name, last_name, phone) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{patient.First_name, patient.Last_name, patient.Phone}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt)
}

func (m PatientModel) Get(id int) (*Patient, error) {
	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT id, created_at, updated_at, first_name, last_name, phone
		FROM patients
		WHERE id = $1
		`
	var patient Patient
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt, &patient.First_name, &patient.Last_name, &patient.Phone)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (m PatientModel) Update(patient *Patient) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE patients
		SET first_name = $1, last_name = $2, phone = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{patient.First_name, patient.Last_name, patient.Phone, patient.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&patient.UpdatedAt)
}

func (m PatientModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM patients
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
