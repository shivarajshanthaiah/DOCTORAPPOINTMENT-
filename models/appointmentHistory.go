package models

type AppointmentHistory struct {
	AppointmentHistoryID int    `gorm:"primaryKey;autoIncrement"`
	AppointmentID        uint   `json:"appointment_id"`
	PatientID            uint   `json:"patient_id"`
	DoctorID             uint   `json:"doctor_id"`
	AppointmentDate      string `json:"appointment_date"`
	AppointmentTime      string `json:"appointment_time"`
	Status               string `json:"status"`
}
