package models

import "github.com/golang-jwt/jwt/v5"

type Doctor struct {
	DoctorID       uint   `gorm:"primaryKey"`
	Name           string `json:"name"`
	Age            string `json:"age"`
	Gender         string `json:"gender"`
	Specialization string `json:"specialization"`
	Experience     string `json:"experience"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Phone          string `json:"phone"`
	LicenseNumber  string `json:"license_number"`
	Availability   bool   `json:"availability"`
	Verified       bool   `json:"verified"`
	HospitalID     uint   `json:"hospital_id"`
}

type DoctorClaims struct {
	Id        uint   `json:"id"`
	DoctorEmail string `json:"useremail"`
	jwt.RegisteredClaims
}

type DoctorAvailability struct {
	DoctorAvailabilityID uint   `gorm:"primaryKey"`
	DoctorID             uint   `json:"doctor_id"`
	Date                 string `json:"date"`
	Available            bool   `json:"available"`
}
