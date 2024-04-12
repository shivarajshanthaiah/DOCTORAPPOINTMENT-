package doctorControllers

import (
	"doctorAppointment/authentication"
	"doctorAppointment/configuration"
	"doctorAppointment/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	DoctorVerified = "false"
	DoctorApproved = "false"
)

func DoctorSignup(c *gin.Context) {
	// Validator instance
	validate := validator.New()

	// Parse doctor data from request
	var doctor models.Doctor
	if err := c.BindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate doctor data
	if err := validate.Struct(doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingDoctor models.Doctor
	if err := configuration.DB.Where("email = ?", doctor.Email).First(&existingDoctor).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	if err := configuration.DB.Where("phone = ?", doctor.Phone).First(&existingDoctor).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
		return
	}

	if err := configuration.DB.Where("license_number = ?", doctor.LicenseNumber).First(&existingDoctor).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Licence already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	doctor.Password = string(hashedPassword)

	// Check if hospital ID exists and is active
	var hospital models.Hospital
	if err := configuration.DB.First(&hospital, doctor.HospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Hospital doesn't exists"})
		return
	}

	if hospital.Status == "Deactive" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Hospital doesn't exists"})
		return
	}
	// Create new doctor record in the database
	doctor.Verified = DoctorVerified
	doctor.Approved = DoctorApproved
	if err := configuration.DB.Create(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reate doctor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor signed up successfully",
		"doctor":  doctor,
	})
}

func DoctorLogin(c *gin.Context) {
	var doctors models.Doctor
	if err := c.BindJSON(&doctors); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	var existingDoctor models.Doctor
	if err := configuration.DB.Where("email = ?", doctors.Email).First(&existingDoctor).Error; err != nil {
		c.JSON(401, gin.H{"error": "invalid is email"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingDoctor.Password), []byte(doctors.Password)); err != nil {
		// Incorrect password
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Check if the doctor is approved
	if existingDoctor.Approved != "true" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor not approved yet"})
		return
	}

	token, err := authentication.GenerateDoctorToken(doctors.Email, doctors.DoctorID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})

}

func ViewHospital(c *gin.Context) {
	var hospitals []models.Hospital

	if err := configuration.DB.Where("status = ?", "Active").Find(&hospitals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hospital not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"Message": "Hospitals list fetehed successfully",
		"data":    hospitals,
	})
}


func DoctorLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You are successfully logged out"})
}


func SaveAvailability(c *gin.Context) {
	var availability models.DoctorAvailability

	if err := c.BindJSON(&availability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if doctor exists
	var doctor models.Doctor
	if err := configuration.DB.Where("doctor_id = ?", availability.DoctorID).First(&doctor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor id not found"})
		return
	}

	if doctor.Approved != "true" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}
	
	var existingAvailability models.DoctorAvailability
    if err := configuration.DB.Where("doctor_id = ? AND date = ?", availability.DoctorID, availability.Date).First(&existingAvailability).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Availability already exists for this date"})
        return
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check availability"})
        return
    }

	if err := configuration.DB.Create(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create availability"})
		return
	}

	c.JSON(http.StatusOK, availability)
}


func AddPrescription(c *gin.Context){
	var prescription models.Prescription
	if err := c.BindJSON(&prescription); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var doctor models.Doctor
	if err := configuration.DB.Where("doctor_id = ?", prescription.DoctorID).First(&doctor).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Invalid doctor ID"})
		return
	}

	var patient models.Patient
	if err := configuration.DB.Where("patient_id = ?", prescription.PatientID).First(&patient).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Invalid patient ID"})
		return
	}

	var appointment models.Appointment
	if err := configuration.DB.Where("appointment_id = ?", prescription.AppointmentID).First(&appointment).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Invalid appointment ID"})
		return
	}

	if err := configuration.DB.Create(&prescription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add prescription"})
		return
	}

	c.JSON(http.StatusOK, prescription)
}