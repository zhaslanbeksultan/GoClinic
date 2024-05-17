package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Appointment struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DateTime  string `json:"date_time"`
	DoctorID  int    `json:"doctor_id"`
	PatientID int    `json:"patient_id"`
}

type AppointmentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m AppointmentModel) Insert(appointment *Appointment) error {
	// Insert a new appointment into the database.
	query := `
		INSERT INTO appointments (date_time, doctor_id, patient_id) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{appointment.DateTime, appointment.DoctorID, appointment.PatientID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt)
}

func (m AppointmentModel) Get(id int) (*Appointment, error) {
	query := `
        SELECT id, created_at, updated_at, doctor_id, patient_id, date_time
        FROM appointments
        WHERE id = $1
    `
	var appointment Appointment
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt, &appointment.DoctorID, &appointment.PatientID, &appointment.DateTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("appointment not found for ID %d", id)
		}
		return nil, err
	}
	return &appointment, nil
}

func (m AppointmentModel) Update(appointment *Appointment) error {
	// Update a specific appointment in the database.
	query := `
		UPDATE appointments
		SET date_time = $1, doctor_id = $2, patient_id = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{appointment.DateTime, appointment.DoctorID, appointment.PatientID, appointment.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&appointment.UpdatedAt)
}

func (m AppointmentModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM appointments
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m AppointmentModel) GetAllSortedByName(filters Filters) ([]*Appointment, error) {
	query := fmt.Sprintf(
		`
       SELECT id, created_at, updated_at, date_time, doctor_id, patient_id
       FROM appointments
       ORDER BY %s %s`,
		filters.sortColumn(),
		filters.sortDirection(),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		if err := rows.Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt, &appointment.DateTime, &appointment.DoctorID, &appointment.PatientID); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (m AppointmentModel) GetFilteredByText(filterText string) ([]*Appointment, error) {
	query := `
       SELECT id, created_at, updated_at, date_time, doctor_id, patient_id
       FROM appointments
       WHERE date_time LIKE '%' || $1 || '%' OR date_time LIKE '%' || $1 || '%'
       ORDER BY date_time
   `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filterText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		if err := rows.Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt, &appointment.DateTime, &appointment.DoctorID, &appointment.PatientID); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (m AppointmentModel) GetPaginatedAppointments(limit, offset int) ([]*Appointment, error) {
	query := `
       SELECT id, created_at, updated_at, date_time, doctor_id, patient_id
       FROM appointments
       ORDER BY id
       LIMIT $1
       OFFSET $2
   `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*Appointment

	for rows.Next() {
		var appointment Appointment
		if err := rows.Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt, &appointment.DateTime, &appointment.DoctorID, &appointment.PatientID); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (m AppointmentModel) Get_By_Doctor(id int) ([]*Appointment, error) {
	query := `
       SELECT id, created_at, updated_at, date_time, doctor_id, patient_id
       FROM appointments
       WHERE doctor_id = $1
       ORDER BY date_time
   `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		if err := rows.Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt, &appointment.DateTime, &appointment.DoctorID, &appointment.PatientID); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (m AppointmentModel) Get_By_Patient(id int) ([]*Appointment, error) {
	query := `
       SELECT id, created_at, updated_at, date_time, doctor_id, patient_id
       FROM appointments
       WHERE patient_id = $1
       ORDER BY date_time
   `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*Appointment
	for rows.Next() {
		var appointment Appointment
		if err := rows.Scan(&appointment.Id, &appointment.CreatedAt, &appointment.UpdatedAt, &appointment.DateTime, &appointment.DoctorID, &appointment.PatientID); err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}
