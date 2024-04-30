package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Doctor struct {
	DoctorID   int    `json:"doctorId"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	Speciality string `json:"speciality"`
	Phone      string `json:"phone"`
}
type DoctorModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *DoctorModel) InsertDoctor(doctor *Doctor) error {
	query := `
			INSERT INTO doctors (firstName, secondName, speciality, phone) 
			VALUES ($1, $2, $3, $4)
			RETURNING doctorId
			`
	args := []interface{}{doctor.FirstName, doctor.SecondName, doctor.Speciality, doctor.Phone}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&doctor.DoctorID)
}
func (m *DoctorModel) GetDoctors(firstName string, secondName string, filters Filters) ([]*Doctor, Metadata, error) {
	query := `
        SELECT COUNT(*) OVER(), d.*
        FROM doctors d
        WHERE (STRPOS(LOWER(d.firstName), LOWER($1)) > 0 OR $1 = '')
        AND (STRPOS(LOWER(d.secondName), LOWER($2)) > 0 OR $2 = '')
        ORDER BY %s %s, d.doctorId ASC
        LIMIT $3 OFFSET $4`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{firstName, secondName, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var doctors []*Doctor
	totalRecords := 0
	for rows.Next() {
		var doctor Doctor
		err := rows.Scan(&totalRecords, &doctor.DoctorID, &doctor.FirstName, &doctor.SecondName, &doctor.Speciality, &doctor.Phone)
		if err != nil {
			return nil, Metadata{}, err
		}
		doctors = append(doctors, &doctor)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return doctors, metadata, nil
}
func (m *DoctorModel) GetDoctor(id int) (*Doctor, error) {
	query := `
        SELECT *
        FROM doctors
        WHERE doctorId = $1
    `
	var doctor Doctor
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&doctor.DoctorID, &doctor.FirstName, &doctor.SecondName, &doctor.Speciality, &doctor.Phone)
	if err != nil {
		return nil, err
	}

	// Теперь запросим список записей для данного доктора
	appointmentsQuery := `
        SELECT *
        FROM appointments
        WHERE doctorId = $1
    `
	rows, err := m.DB.QueryContext(ctx, appointmentsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return &doctor, nil
}
func (m *DoctorModel) UpdateDoctor(doctor *Doctor) error {
	query := `
			UPDATE doctors
			SET firstName = $1,
			secondName = $2,
			speciality = $3,
			phone = $4
			WHERE doctorId = $5
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, doctor.FirstName, doctor.Speciality, doctor.Speciality, doctor.Phone, doctor.DoctorID)
	return err
}
func (m *DoctorModel) DeleteDoctor(id int) error {
	query := `
			DELETE FROM doctors
			WHERE doctorId = $1
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
