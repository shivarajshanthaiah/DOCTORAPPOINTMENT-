package models


type Appointment struct {
	AppointmentID      int    `json:"appointment_id"`
	PatientID          int    `json:"patient_id"`
	DoctorID           int    `json:"doctor_id"`
	AppointmentDate    string `json:"appointment_date"`
	AppointmentTime    string `json:"appointment_time"`
	Status             string `json:"status"`
	PatientHealthIssue string `json:"patient_health_issue"`
}