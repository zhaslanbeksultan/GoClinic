package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Patient struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type PatientModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m PatientModel) Insert(patient *Patient) error {
	// Insert a new patient into the database.
	query := `
		INSERT INTO patients (first_name, last_name, phone) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{patient.FirstName, patient.LastName, patient.Phone}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt)
}

func (m PatientModel) Get(id int) (*Patient, error) {
	// Retrieve a specific patient based on his ID.
	query := `
		SELECT id, created_at, updated_at, first_name, last_name, phone
		FROM patients
		WHERE id = $1
		`
	var patient Patient
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt, &patient.FirstName, &patient.LastName, &patient.Phone)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (m PatientModel) Update(patient *Patient) error {
	// Update a specific patient in the database.
	query := `
		UPDATE patients
		SET first_name = $1, last_name = $2, phone = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{patient.FirstName, patient.LastName, patient.Phone, patient.Id}
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

func (m PatientModel) GetAllSortedByName() ([]*Patient, error) {
	query := `
        SELECT id, created_at, updated_at, first_name, last_name, phone
        FROM patients
        ORDER BY first_name, last_name
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
		if err := rows.Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt, &patient.FirstName, &patient.LastName, &patient.Phone); err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

func (m PatientModel) GetFilteredByText(filterText string) ([]*Patient, error) {
	// Определение запроса для фильтрации пациентов
	query := `
        SELECT id, created_at, updated_at, first_name, last_name, phone
        FROM patients
        WHERE first_name LIKE '%' || $1 || '%' OR last_name LIKE '%' || $1 || '%'
        ORDER BY first_name, last_name
    `

	// Выполнение запроса к базе данных с параметром фильтрации и получение результатов
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filterText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*Patient
	for rows.Next() {
		var patient Patient
		if err := rows.Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt, &patient.FirstName, &patient.LastName, &patient.Phone); err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

func (m PatientModel) GetPatients(limit, offset int) ([]*Patient, error) {
	// Определение запроса для получения пациентов с учетом лимита и смещения
	query := `
        SELECT id, created_at, updated_at, first_name, last_name, phone
        FROM patients
        ORDER BY id
        LIMIT $1
        OFFSET $2
    `

	// Выполнение запроса к базе данных с использованием лимита и смещения
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Инициализация массива для хранения результатов
	var patients []*Patient

	// Парсинг результатов запроса и добавление их в массив
	for rows.Next() {
		var patient Patient
		if err := rows.Scan(&patient.Id, &patient.CreatedAt, &patient.UpdatedAt, &patient.FirstName, &patient.LastName, &patient.Phone); err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}

	// Проверка наличия ошибок после парсинга
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возвращение массива пациентов
	return patients, nil
}

func GetPaginatedPatients(limit, offset int, m PatientModel) ([]*Patient, error) {
	if limit < 1 || offset < 0 {
		return nil, errors.New("Invalid limit or offset parameter")
	}

	// Получение пациентов из базы данных с использованием переданного лимита и смещения
	patients, err := m.GetPatients(limit, offset)
	if err != nil {
		return nil, err
	}
	return patients, nil
}
