package models

import "time"

type Slot struct {
	SlotID    int       `gorm:"primaryKey;autoIncremet"`
	DoctorID  int       `json:"doctor_id"`
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Available bool      `json:"available"`
}
