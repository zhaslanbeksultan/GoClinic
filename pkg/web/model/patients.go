package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Patient struct {
	PatientID int    `json:"patientId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}
type PatientModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *PatientModel) InsertPatient(patient *Patient) error {
	query := `
			INSERT INTO patients (firstName, lastName, phone)
			VALUES ($1, $2, $3)
			RETURNING patientId
			`
	args := []interface{}{patient.FirstName, patient.LastName, patient.Phone}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&patient.PatientID)
}
func (m *PatientModel) GetPatient(id int) (*Patient, error) {
	query := `
			SELECT patientId, firstName, lastName, phone
			FROM patients
			WHERE patientId = $1
			`
	var patient Patient
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&patient.PatientID, &patient.FirstName, &patient.LastName, &patient.Phone)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}
func (m *PatientModel) GetPatients() ([]*Patient, error) {
	query := `
			SELECT patientId, firstName, lastName, phone
			FROM patients
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var patients []*Patient
	for rows.Next() {
		var patient Patient
		err := rows.Scan(&patient.PatientID, &patient.FirstName, &patient.LastName, &patient.Phone)
		if err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return patients, nil
}
func (m *PatientModel) UpdatePatient(patient *Patient) error {
	query := `
			UPDATE patients
			SET firstName = $1, lastName = $2, phone = $3
			WHERE patientId = $8
			`
	args := []interface{}{patient.FirstName, patient.LastName, patient.Phone, patient.PatientID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
func (m *PatientModel) DeletePatient(id int) error {
	query := `
			DELETE FROM patients
			WHERE patientId = $1
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
