package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Doctor struct {
	Id         string `json:"id"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Speciality string `json:"speciality"`
	Phone      string `json:"phone"`
}

type DoctorModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m DoctorModel) Insert(doctor *Doctor) error {
	// Insert a new doctor into the database.
	query := `
		INSERT INTO doctors (first_name, last_name, speciality, phone) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{doctor.FirstName, doctor.LastName, doctor.Speciality, doctor.Phone}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&doctor.Id, &doctor.CreatedAt, &doctor.UpdatedAt)
}

func (m DoctorModel) Get(id int) (*Doctor, error) {
	// Retrieve a specific Doctor based on his ID.
	query := `
		SELECT id, created_at, updated_at, first_name, last_name, speciality, phone
		FROM doctors
		WHERE id = $1
		`
	var doctor Doctor
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&doctor.Id, &doctor.CreatedAt, &doctor.UpdatedAt, &doctor.FirstName, &doctor.LastName, &doctor.Speciality, &doctor.Phone)
	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (m DoctorModel) Update(doctor *Doctor) error {
	// Update a specific doctor in the database.
	query := `
		UPDATE doctors
		SET first_name = $1, last_name = $2, speciality = $3, phone = $4
		WHERE id = $5
		RETURNING updated_at
		`
	args := []interface{}{doctor.FirstName, doctor.LastName, doctor.Speciality, doctor.Phone, doctor.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&doctor.UpdatedAt)
}

func (m DoctorModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM doctors
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m DoctorModel) GetAllSortedByName(filters Filters) ([]*Doctor, error) {
	query := fmt.Sprintf(
		`
        SELECT id, created_at, updated_at, first_name, last_name, speciality, phone
        FROM doctors
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

	var doctors []*Doctor
	for rows.Next() {
		var doctor Doctor
		if err := rows.Scan(&doctor.Id, &doctor.CreatedAt, &doctor.UpdatedAt, &doctor.FirstName, &doctor.LastName, &doctor.Speciality, &doctor.Phone); err != nil {
			return nil, err
		}
		doctors = append(doctors, &doctor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (m DoctorModel) GetFilteredByText(filterText string) ([]*Doctor, error) {
	query := `
        SELECT id, created_at, updated_at, first_name, last_name, speciality, phone
        FROM doctors
        WHERE first_name LIKE '%' || $1 || '%' OR last_name LIKE '%' || $1 || '%'
        ORDER BY first_name, last_name
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filterText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []*Doctor
	for rows.Next() {
		var doctor Doctor
		if err := rows.Scan(&doctor.Id, &doctor.CreatedAt, &doctor.UpdatedAt, &doctor.FirstName, &doctor.LastName, &doctor.Speciality, &doctor.Phone); err != nil {
			return nil, err
		}
		doctors = append(doctors, &doctor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (m DoctorModel) GetPaginatedDoctors(limit, offset int) ([]*Doctor, error) {
	query := `
        SELECT id, created_at, updated_at, first_name, last_name, speciality, phone
        FROM doctors
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

	var doctors []*Doctor

	for rows.Next() {
		var doctor Doctor
		if err := rows.Scan(&doctor.Id, &doctor.CreatedAt, &doctor.UpdatedAt, &doctor.FirstName, &doctor.LastName, &doctor.Speciality, &doctor.Phone); err != nil {
			return nil, err
		}
		doctors = append(doctors, &doctor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}
