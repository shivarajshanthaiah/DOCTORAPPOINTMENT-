package models

type Hospital struct {
    HospitalID  int     `gorm:"primaryKey;autoIncrement"`
    Name        string  `json:"name"`
    Location    string  `json:"location"`
}