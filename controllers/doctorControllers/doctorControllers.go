package doctorControllers

import (
	"doctorAppointment/configuration"
	"doctorAppointment/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)




// ViewHospital retrieves a list of active hospitals
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

// DoctorLogout
func DoctorLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You are successfully logged out"})
}

// SaveAvailability saves the availability of a doctor
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

	// Check if doctor is approved
	if doctor.Approved != "true" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}
	
	// Check if availability for the given date already exists
	var existingAvailability models.DoctorAvailability
    if err := configuration.DB.Where("doctor_id = ? AND date = ?", availability.DoctorID, availability.Date).First(&existingAvailability).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Availability already exists for this date"})
        return
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check availability"})
        return
    }

	// Create new availability record in the database
	if err := configuration.DB.Create(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create availability"})
		return
	}

	c.JSON(http.StatusOK, availability)
}


// AddPrescription
func AddPrescription(c *gin.Context){
	var prescription models.Prescription
	if err := c.BindJSON(&prescription); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if doctor exists
	var doctor models.Doctor
	if err := configuration.DB.Where("doctor_id = ?", prescription.DoctorID).First(&doctor).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Invalid doctor ID"})
		return
	}

	// Check if patient exists
	var patient models.Patient
	if err := configuration.DB.Where("patient_id = ?", prescription.PatientID).First(&patient).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Invalid patient ID"})
		return
	}

	// Check if appointment exists
	var appointment models.Appointment
	if err := configuration.DB.Where("appointment_id = ?", prescription.AppointmentID).First(&appointment).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Invalid appointment ID"})
		return
	}

	// Create new prescription record in the database
	if err := configuration.DB.Create(&prescription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add prescription"})
		return
	}

	c.JSON(http.StatusOK, prescription)
}