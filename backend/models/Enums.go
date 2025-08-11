package models

type Role string

const (
	RoleAdmin        Role = "ADMIN"
	RolePatient      Role = "PATIENT"
	RoleDoctor       Role = "DOCTOR"
	RoleReceptionist Role = "RECEPTIONIST"
)

type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "PENDING"
	AppointmentStatusConfirmed AppointmentStatus = "CONFIRMED"
	AppointmentStatusCancelled AppointmentStatus = "CANCELLED"
	AppointmentStatusCompleted AppointmentStatus = "COMPLETED"
)

type VitalType string

const (
	BP     VitalType = "BLOOD_PRESSURE"
	HR     VitalType = "HEART_RATE"
	Weight VitalType = "WEIGHT"
	BMI    VitalType = "BMI"
	Temp   VitalType = "TEMPERATURE"
	RR     VitalType = "RESPIRATORY_RATE"
	OS     VitalType = "OXYGEN_SATURATION"
)

type GoalType string

const (
	Steps GoalType = "STEPS"
	Water GoalType = "WATER"
	Sleep GoalType = "SLEEP"
)

type ApptType string

const (
	ApptTypeConsultation ApptType = "CONSULTATION"
	ApptTypeFollowup     ApptType = "FOLLOWUP"
	ApptTypeCheckup      ApptType = "CHECKUP"
	ApptTypeEmergency    ApptType = "EMERGENCY"
)
