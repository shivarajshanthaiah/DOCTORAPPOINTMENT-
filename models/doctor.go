package models

import "github.com/golang-jwt/jwt/v5"

type Doctor struct {
	DoctorID       int    `gorm:"primaryKey"`
	Name           string `json:"name"`
	Age            int    `json:"age"`
	Gender         string `json:"gender"`
	Specialization string `json:"specialization"`
	Experience     string `json:"experience"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Phone          string `json:"phone"`
	LicenseNumber  string `json:"license_number"`
	Availability   bool   `json:"availability"`
	Verified       bool   `json:"verified"`
	HospitalID     int    `json:"hospital_id"`
}

type DoctorClaims struct {
	Id        uint   `json:"id"`
	UserEmail string `json:"useremail"`
	jwt.RegisteredClaims
}

type DoctorAvailability struct {
	DoctorAvailabilityID int    `gorm:"primaryKey"`
	DoctorID             int    `json:"doctor_id"`
	Date                 string `json:"date"`
	Available            bool   `json:"available"`
}
