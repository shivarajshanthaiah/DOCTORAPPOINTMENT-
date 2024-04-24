package adminControllers

import (
	"doctorAppointment/configuration"
	"doctorAppointment/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBookingStatusCounts(c *gin.Context) {
	var totalBookings int64
	result := configuration.DB.Model(&models.Appointment{}).Count(&totalBookings)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch total bookings"})
		return
	}

	var confirmedBookings int64
	confirmedResults := configuration.DB.Model(&models.Appointment{}).Where("booking_status = ?", "confirmed").Count(&confirmedBookings)
	if confirmedResults.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch confirmed bookings"})
		return
	}

	var completedBookings int64
	completedResult := configuration.DB.Model(&models.Appointment{}).Where("booking_status = ?", "completed").Count(&completedBookings)
	if completedResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch completed bookings"})
		return
	}

	var cancelledBookings int64
	cancelledResult := configuration.DB.Model(&models.Appointment{}).Where("booking_status = ?", "cancelled").Count(&cancelledBookings)
	if cancelledResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch cancelled bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":            "Sucess",
		"Message":           "Booking details fetched sucessfully",
		"TotalBookings":     totalBookings,
		"ConfirmedBookings": confirmedBookings,
		"CompletedBookings": completedBookings,
		"CancelledBookings": cancelledBookings,
	})
}

type DoctorBooking struct {
	DoctorID     int `jsoon:"doctor_id"`
	BookingCount int `json:"booking_count"`
}

func GetDoctorWiseBookings(c *gin.Context) {
	var doctorBookings []DoctorBooking
	result := configuration.DB.Model(&models.Appointment{}).
		Select("doctor_id, COUNT(*) as booking_count").
		Group("doctor_id").
		Find(&doctorBookings)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch doctor-wise bookings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status":         "Success",
		"Message":        "Doctor-wise bookings fetched successfully",
		"DoctorBookings": doctorBookings,
	})
}

type DepartmentBooking struct {
	Specialization string `json:"specialization"`
	BookingCount   int    `json:"booking_count"`
}

func GetDepartmentWiseBookings(c *gin.Context) {
	var departmentBookings []DepartmentBooking

	result := configuration.DB.Model(&models.Appointment{}).
		Select("doctors.specialization as specialization, COUNT(*) as booking_count").
		Joins("JOIN doctors ON appointments.doctor_id = doctors.doctor_id").
		Group("doctors.specialization").
		Find(&departmentBookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch department-wise booking details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":             "Success",
		"Message":            "Details fetched succesfully",
		"DepartmentBookings": departmentBookings,
	})
}

type Revenue struct {
	Day   *float64 `json:"day"`
	Week  *float64 `json:"week"`
	Month *float64 `json:"month"`
	Year  *float64 `json:"year"`
}

func GetTotalRevenue(c *gin.Context) {
	now := time.Now()

	// Get the start and end time for the day
	startofDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endofDay := startofDay.AddDate(0, 0, 1).Add(-time.Second)

	// Get the start and end time for the week
	startofWeek := startofDay.AddDate(0, 0, -int(now.Weekday()))
	endofWeek := startofWeek.AddDate(0, 0, 7).Add(-time.Second)

	// Get the start and end time for the month
	startofMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endofMonth := startofMonth.AddDate(0, 1, 0).Add(-time.Second)

	// Get the start and end time for the year
	startofYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	endofYear := startofYear.AddDate(1, 0, 0).Add(-time.Second)

	// Query the database to get the total revenue for different timeframes
	var revenue Revenue
	result := configuration.DB.Model(&models.Invoice{}).
		Select("SUM(total_amount) as total_revenue").
		Where("payment_status = ?", "Paid").
		Where("updated_at BETWEEN ? AND ?", startofDay, endofDay).
		Scan(&revenue.Day)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch revenue for the day"})
		return
	}

	result = configuration.DB.Model(&models.Invoice{}).
		Select("SUM(total_amount) as total_revenue").
		Where("payment_status = ?", "Paid").
		Where("updated_at BETWEEN ? AND ?", startofWeek, endofWeek).
		Scan(&revenue.Week)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch revenue for the week"})
		return
	}

	result = configuration.DB.Model(&models.Invoice{}).
		Select("SUM(total_amount) as total_revenue").
		Where("payment_status = ?", "Paid").
		Where("updated_at BETWEEN ? AND ?", startofMonth, endofMonth).
		Scan(&revenue.Month)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch revenue for the month"})
		return
	}

	result = configuration.DB.Model(&models.Invoice{}).
		Select("SUM(total_amount) as total_revenue").
		Where("payment_status = ?", "Paid").
		Where("updated_at BETWEEN ? AND ?", startofYear, endofYear).
		Scan(&revenue.Year)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch revenue for the year"})
		return
	}

	// Respond with the total revenue for different timeframes
	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"message": "Revenure details fetched sucessfully",
		"Revenue": revenue,
	})
}

type SpecificRevenue struct {
	Revenue *float64 `json:"revenue"`
}

func GetSpecificRevenue(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid start date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		startDate = time.Now()
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid end date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		endDate = time.Now()
	}

	var specificRevenue SpecificRevenue
	result := configuration.DB.Model(&models.Invoice{}).
		Select("SUM(total_amount) as total_revenue").
		Where("payment_status = ?", "Paid").
		Where("updated_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&specificRevenue.Revenue)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch revenue for specific date range"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"Message": "Revenue details fetched successfully",
		"Revenue": specificRevenue,
	})
}
