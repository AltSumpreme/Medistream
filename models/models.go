package models

func AllModels() []interface{} {
	return []interface{}{
		&Auth{},
		&User{},
		&Appointment{},
		&Doctor{},
		&Patient{},
		&RefreshToken{},
		&MedicalRecord{},
		&Report{},
		&Vital{},
		&Message{},
		&Prescription{},
	}
}
