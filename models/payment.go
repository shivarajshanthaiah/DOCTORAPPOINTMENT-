package models

import (

	"gorm.io/gorm"
)

type Payment struct {
    gorm.Model
	AppointmentID uint      `json:"appointment_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
}