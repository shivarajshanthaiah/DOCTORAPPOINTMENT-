package controllers

import (
	"doctorAppointment/configuration"
	"doctorAppointment/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserWallet helps to get user wallet by user id
func Wallet(c *gin.Context) {
	userid := c.Param("userid")

	var wallet models.Wallet
	if err := configuration.DB.Where("user_id = ?", userid).First(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to find user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status":        "Success",
		"Wallet Amount": wallet.Amount,
	})

}

// CancelAppointment is a handler function for cancelling an appointment.
// func CancelAppointment(c *gin.Context) {
// 	var appointment models.Appointment
// 	if err := configuration.DB.Where("appointment_id = ?", c.Param("id")).First(&appointment).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
// 		return
// 	}

// 	if appointment.BookingStatus == "cancelled"{
// 		c.JSON(http.StatusBadRequest, gin.H{"Error":"Appointments has been cancelled already"})
// 		return
// 	}

// 	if appointment.BookingStatus == "completed"{
// 		c.JSON(http.StatusBadRequest, gin.H{"Error":"This appointment has already been completed"})
// 		return
// 	}

// 	if appointment.BookingStatus != "confirmed" {
// 		c.JSON(http.StatusBadRequest, gin.H{"Error": "Appointmenet cannot be cancelled as it is not confirmed"})
// 		return
// 	}

// 	var invoice models.Invoice
// 	if err := configuration.DB.Where("appointment_id = ?", c.Param("id")).First(&invoice).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
// 		return
// 	}
// 	if invoice.PaymentMethod != "online" {
// 		c.JSON(http.StatusBadRequest, gin.H{"Error": "Appointmenet payment was not made online, refund not applicable"})
// 		return
// 	}

// 	refundAmount := invoice.TotalAmount * 0.95

// 	var wallet models.Wallet
// 	if err := configuration.DB.Where("user_id = ?", uint(appointment.PatientID)).First(&wallet).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// Wallet doesn't exist, create a new one
// 			wallet = models.Wallet{
// 				UserID: uint(appointment.PatientID),
// 				Amount: refundAmount,
// 			}
// 			if err := configuration.DB.Create(&wallet).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create wallet"})
// 				return
// 			}
// 		} else {
// 			// Error occurred while fetching wallet
// 			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch wallet"})
// 			return
// 		}
// 	} else {
// 		// Wallet exists, update its amount
// 		wallet.Amount += refundAmount
// 		if err := configuration.DB.Save(&wallet).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update wallet"})
// 			return
// 		}
// 	}

// 	appointment.BookingStatus = "cancelled"
// 	if err := configuration.DB.Save(&appointment).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update appointment status"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": fmt.Sprintf("Appointment Cancelled. Refund amount : %.2f", refundAmount),
// 	})
// }

// CancelAppointment is a handler function for cancelling an appointment.
func CancelAppointment(c *gin.Context) {
	var appointment models.Appointment
	if err := configuration.DB.Where("appointment_id = ?", c.Param("id")).First(&appointment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	if appointment.BookingStatus == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Appointment has already been cancelled"})
		return
	}

	if appointment.BookingStatus == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "This appointment has already been completed"})
		return
	}

	if appointment.BookingStatus != "confirmed" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Appointment cannot be cancelled as it is not confirmed"})
		return
	}

	var invoice models.Invoice
	if err := configuration.DB.Where("appointment_id = ?", c.Param("id")).First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	if invoice.PaymentMethod == "online" {
		// Refund applicable for online payments
		refundAmount := invoice.TotalAmount * 0.95

		// Update payment status to refunded
		invoice.PaymentStatus = "refunded"

		// Update invoice in the database
		if err := configuration.DB.Save(&invoice).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update invoice"})
			return
		}

		var wallet models.Wallet
		if err := configuration.DB.Where("user_id = ?", appointment.PatientID).First(&wallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Wallet doesn't exist, create a new one
				wallet = models.Wallet{
					UserID: uint(appointment.PatientID),
					Amount: refundAmount,
				}
				if err := configuration.DB.Create(&wallet).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create wallet"})
					return
				}
			} else {
				// Error occurred while fetching wallet
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch wallet"})
				return
			}
		} else {
			// Wallet exists, update its amount
			wallet.Amount += refundAmount
			if err := configuration.DB.Save(&wallet).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update wallet"})
				return
			}
		}

		appointment.BookingStatus = "cancelled"
		if err := configuration.DB.Save(&appointment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update appointment status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Appointment Cancelled. Refund amount: %.2f", refundAmount),
		})
	} else {
		// For offline payments, simply cancel the appointment without refunding
		appointment.BookingStatus = "cancelled"

		// Update appointment status in the database
		if err := configuration.DB.Save(&appointment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update appointment status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Appointment Cancelled. Amount cannot be refunded as payment method was not online"})
	}
}
