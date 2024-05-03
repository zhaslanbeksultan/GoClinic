package filler

import (
	model "github.com/zhaslanbeksultan/GoClinic/pkg/web/model"
)

func PopulateDatabase(models model.Models) error {
	for _, doctor := range doctors {
		models.Doctors.InsertDoctor(&doctor)
	}
	// TODO: Implement restaurants pupulation
	// TODO: Implement the relationship between restaurants and menus
	return nil
}

var doctors = []model.Doctor{
	{FirstName: "Caesar", SecondName: "Salad", Speciality: "Nutritionist", Phone: "123-456-7890"},
	{FirstName: "Greek", SecondName: "Salad", Speciality: "Dietitian", Phone: "234-567-8901"},
	{FirstName: "Caresse", SecondName: "Salad", Speciality: "Health Coach", Phone: "345-678-9012"},
	{FirstName: "Cobb", SecondName: "Salad", Speciality: "Fitness Trainer", Phone: "456-789-0123"},
	{FirstName: "Kale", SecondName: "Salad", Speciality: "Wellness Consultant", Phone: "567-890-1234"},
}
