package model

import (
	"database/sql"
	"errors"
	"log"
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

var doctors = []Doctor{
	{
		Id:         "1",
		FirstName:  "Azat",
		LastName:   "Azatov",
		Speciality: "Neurosurgeon",
		Phone:      "87770820010",
	},
	{
		Id:         "2",
		FirstName:  "Beks",
		LastName:   "Zhaslan",
		Speciality: "Dentist",
		Phone:      "87778250062",
	},
	{
		Id:         "3",
		FirstName:  "Aidar",
		LastName:   "Zhaxylykuly",
		Speciality: "Neurologist",
		Phone:      "87788522332",
	}, {
		Id:         "4",
		FirstName:  "Dada",
		LastName:   "Didov",
		Speciality: "Therapist",
		Phone:      "87478256639",
	}, {
		Id:         "5",
		FirstName:  "Cruc",
		LastName:   "Cece",
		Speciality: "Pediatrician",
		Phone:      "87085056640",
	},
}

func GetDoctors() []Doctor {
	return doctors
}

func GetDoctor(id string) (*Doctor, error) {
	for _, r := range doctors {
		if r.Id == id {
			return &r, nil
		}
	}
	return nil, errors.New("doctor not found")
}
