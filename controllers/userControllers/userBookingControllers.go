package userControllers

import (
	"doctorAppointment/configuration"
	"doctorAppointment/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAvailableTimeSlots(c *gin.Context) {
	doctorID := c.Param("doctor_id")
	dateStr := c.Query("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	fmt.Println("im here")
	var availability models.DoctorAvailability
	if err := configuration.DB.Where("doctor_id = ? AND date = ?", doctorID, date).First(&availability).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Availability not found"})
		return
	}
	fmt.Println("not here")

	startTime, endTime := splitAvailabilityTime(availability.AvilableTime)
	availableTimeSlots := divideSlots(startTime, endTime, 30*time.Minute)

	var bookings []models.Appointment
	if err := configuration.DB.Where("doctor_id = ? AND appointment_date = ?", doctorID, date).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}

	fmt.Print("not exaxtly")

	bookedTimeSlots := make(map[string]bool)
	for _, booking := range bookings {
		bookedTimeSlots[booking.AppointmentTimeSlot] = true
	}

	adjustedTimeSlots := make([]string, 0)
	for _, slot := range availableTimeSlots {
		if !bookedTimeSlots[slot] {
			adjustedTimeSlots = append(adjustedTimeSlots, slot)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"date":                 dateStr,
		"available_time_slots": adjustedTimeSlots,
	})
}

func splitAvailabilityTime(availabilityTime string) (startTime, endTime string) {
	parts := strings.Split(availabilityTime, "-")
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func divideSlots(startTime, endTime string, interval time.Duration) []string {

	start, _ := time.Parse("15:04", startTime)
	end, _ := time.Parse("15:04", endTime)

	fmt.Println(start)
	var slots []string
	for t := start; t.Before(end); t = t.Add(interval) {
		slotEnd := t.Add(interval)
		slots = append(slots, fmt.Sprintf("%s-%s", t.Format("15:04"), slotEnd.Format("15:04")))
	}
	return slots
}

type DoctorInfo struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Gender     string `json:"gender" gorm:"not null"`
	Speciality string `json:"speciality"`
	Experience int    `json:"experience"`
	Location string
}

func GetDoctorsBySpeciality(c *gin.Context) {
	var doctors []models.Doctor
	doctorSpeciality := c.Param("specialization")
	if err := configuration.DB.Where("specialization = ? AND approved = ?", doctorSpeciality, "true").Find(&doctors).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No doctors found with the specified speciality"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   "Couldn't Get doctors details",
			"details": err.Error()})
		return
	}

	var doctorInfoList []DoctorInfo
	for _, doctor := range doctors {
		var hospital models.Hospital
		err := configuration.DB.Where("id = ?",doctor.HospitalID).First(&hospital).Error
		if err != nil{
			c.JSON(http.StatusNotFound, gin.H{"error":"Location error"})
			return
		}
		doctorInfo := DoctorInfo{
			Name:       doctor.Name,
			Age:        doctor.Age,
			Gender:     doctor.Gender,
			Speciality: doctor.Specialization,
			Experience: doctor.Experience,
			Location: hospital.Location,
		}
		doctorInfoList = append(doctorInfoList, doctorInfo)
	}

	if len(doctors) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No doctors found with the specified speciality"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"Message": "Doctors details list fetched successfully",
		"data":    doctorInfoList,
	})
}

// func BookAppointment(c *gin.Context) {
// 	var booking models.Appointment

// 	if err := c.BindJSON(&booking); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	existingAppointment := models.Appointment{}
// 	err := configuration.DB.Where("doctor_id = ? AND appointment_date = ? AND appointment_time_slot = ?", booking.DoctorID, booking.AppointmentDate, booking.AppointmentTimeSlot).First(&existingAppointment).Error
// 	if err == nil {
// 		c.JSON(http.StatusConflict, gin.H{"error": "Appointment already booked for the same doctor, date and time slot"})
// 		return
// 	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing appointments"})
// 		return
// 	}

// 	var doctor models.Doctor
// 	if err := configuration.DB.Where("doctor_id = ? AND approved = ?", booking.DoctorID, "true").First(&doctor).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"Error": "Doctor not found"})
// 		return
// 	}
// 	 // Check if the patient exists
// 	 var patient models.Patient
// 	 if err := configuration.DB.Where("patient_id = ?", booking.PatientID).First(&patient).Error; err != nil {
// 		 c.JSON(http.StatusNotFound, gin.H{"error": "Wrong patient ID"})
// 		 return
// 	 }

// 	if err := configuration.DB.Create(&booking).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book appointment"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, booking)
// }


func BookAppointment(c *gin.Context) {
	var booking models.Appointment

	if err := c.BindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the appointment time slot is within the available time slots of the doctor
	doctorAvailability := getDoctorAvailability(booking.DoctorID, booking.AppointmentDate)
	if doctorAvailability == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor availability not found"})
		return
	}

	availableTimeSlots := divideAvailableSlots(doctorAvailability.AvilableTime, 30*time.Minute)

	if !isTimeWithinAvailableSlot(booking.AppointmentTimeSlot, availableTimeSlots) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment time slot not available"})
		return
	}

	// Check for existing appointments with the same date and time slot
	if !isAppointmentAvailable(booking.DoctorID, booking.AppointmentDate, booking.AppointmentTimeSlot,) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment already booked for the same date and time slot wiht the doctor"})
		return
	}

	// Check if the patient exists
	var patient models.Patient
	if err := configuration.DB.Where("patient_id = ?", booking.PatientID).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wrong patient ID"})
		return
	}

	if !isDuplicateAppointment(booking.DoctorID, booking.AppointmentDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment already booked with the same doctor in the same day"})
		return
	}
	// Create the appointment
	if err := configuration.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book appointment"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func getDoctorAvailability(doctorID int, date time.Time) *models.DoctorAvailability {
	var availability models.DoctorAvailability
	if err := configuration.DB.Where("doctor_id = ? AND date = ?", doctorID, date).First(&availability).Error; err != nil {
		return nil
	}
	return &availability
}

func isTimeWithinAvailableSlot(appointmentTimeSlot string, availableSlots []string) bool {
	for _, slot := range availableSlots {
		if slot == appointmentTimeSlot {
			return true
		}
	}
	return false
}

func isAppointmentAvailable(doctorID int, date time.Time, appointmentTimeSlot string,) bool {
	var existingAppointment models.Appointment
	err := configuration.DB.Where("doctor_id = ? AND appointment_date = ? AND appointment_time_slot = ?", doctorID, date, appointmentTimeSlot,).First(&existingAppointment).Error
	if err == nil {
		return false // Appointment already exists for the same date and time slot
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// An unexpected error occurred while querying the database
		log.Println("Error checking for existing appointment:", err)
		return false
	}
	return true
}

func divideAvailableSlots(availability string, interval time.Duration) []string {
	// Extract start and end times from the availability string
	parts := strings.Split(availability, "-")
	if len(parts) != 2 {
		fmt.Println("Invalid availability format")
		return nil
	}
	startTime := parts[0]
	endTime := parts[1]

	// Parse start and end times
	start, _ := time.Parse("15:04", startTime)
	end, _ := time.Parse("15:04", endTime)

	var slots []string
	for t := start; t.Before(end); t = t.Add(interval) {
		slotEnd := t.Add(interval)
		slots = append(slots, fmt.Sprintf("%s-%s", t.Format("15:04"), slotEnd.Format("15:04")))
	}
	return slots
}

func isDuplicateAppointment(doctorID int, date time.Time) bool {
	var existingAppointments []models.Appointment
	err := configuration.DB.Where("doctor_id = ? AND appointment_date = ?", doctorID, date).Find(&existingAppointments).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true // No existing appointments found for the same doctor and date
		}
		// An unexpected error occurred while querying the database
		log.Println("Error checking for existing appointments:", err)
		return false
	}
	// Found existing appointments for the same doctor and date
	return len(existingAppointments) == 0
}