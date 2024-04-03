package models

import (

	"gorm.io/gorm"
)

type Prescription struct {
	gorm.Model
	DoctorID         int       `json:"doctor_id"`
	PatientID        int       `json:"patient_id"`
	AppointmentID    int       `json:"appointment_id"`
	PrescriptionText string    `json:"prescription_text"`
}