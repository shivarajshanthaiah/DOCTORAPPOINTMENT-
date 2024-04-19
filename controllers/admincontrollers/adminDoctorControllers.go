package adminControllers

import (
	"doctorAppointment/configuration"
	"doctorAppointment/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// View verified doctors
func ViewVerifiedDoctors(c *gin.Context) {
	var doctors []models.Doctor

	if err := configuration.DB.Where("verified = ?", "true").Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching verified Doctors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"Message": "Verified doctors list fetched sucessfully",
		"data":    doctors,
	})
}

//View Not verified doctors
func ViewNotVerifiedDoctors(c *gin.Context) {
	var doctors []models.Doctor

	if err := configuration.DB.Where("verified = ?", "false").Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching Doctors list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"Message": "Doctors list fetched sucessfully",
		"data":    doctors,
	})
}

// View Verified and approved dospitals
func ViewVerifiedApprovedDoctors(c *gin.Context) {
    var doctors []models.Doctor

    if err := configuration.DB.Where("verified = ? AND approved = ?", "true", "true").Find(&doctors).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching verified and approved Doctors"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "Status":  "Success",
        "Message": "Verified and approved doctors list fetched successfully",
        "data":    doctors,
    })
}

//View verified but not approved doctors
func ViewVerifiedNotApprovedDoctors(c *gin.Context) {
    var doctors []models.Doctor

    if err := configuration.DB.Where("verified = ? AND approved = ?", "true", "false").Find(&doctors).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching verified and approved Doctors"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "Status":  "Success",
        "Message": "Verified and approved doctors list fetched successfully",
        "data":    doctors,
    })
}

//Update doctor credentials
func UpdateDoctor(c *gin.Context) {
	var doctor models.Doctor
	doctorID := c.Param("id")

	if err := configuration.DB.First(&doctor, doctorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "No doctor with this ID"})
		return
	}

	if err := c.BindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	if err := configuration.DB.Save(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"message": "Doctor detailes have been updated sucessfully sucessfully",
		"data":    doctor,
	})
}

//View all Doctors list
func ViewDoctors(c *gin.Context) {
	var doctors []models.Doctor

	if err := configuration.DB.Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Doctors not found"})
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"Message": "Doctors list fetehed successfully",
		"data":    doctors,
	})
}

//Get doctors details by id
func GetDoctorByID(c *gin.Context) {
	var doctor models.Doctor
	doctorID := c.Param("id")

	if err := configuration.DB.First(&doctor, doctorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Couldn't Get doctor details"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"Message": "Doctor details fetched successfully",
		"data":    doctor,
	})

}

//Get doctros details by speciality
func GetDoctorBySpeciality(c *gin.Context) {
	var doctors []models.Doctor
	doctorSpeciality := c.Param("specialization")
	if err := configuration.DB.Where("specialization = ?", doctorSpeciality).Find(&doctors).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No doctors found with the specified speciality"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   "Couldn't Get doctors details",
			"details": err.Error()})
		return
	}
	if len(doctors) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No doctors found with the specified speciality"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"Message": "Doctors details list fetched successfully",
		"data":    doctors,
	})
}
