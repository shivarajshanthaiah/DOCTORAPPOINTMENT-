package controllers

import (
	"doctorAppointment/configuration"
	"doctorAppointment/models"
	"errors"
	"fmt"

	// "fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

// func GetInvoice(c *gin.Context) {
// 	var invoice []models.Invoice
// 	if err := configuration.DB.Find(&invoice).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error occured while receiving the invoice",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, invoice)
// }

// To make payment offline
func PayInvoiceOffline(c *gin.Context) {
	var paymentRequest struct {
		InvoiceID uint `json:"invoice_id"`
	}

	if err := c.BindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var invoice models.Invoice
	if err := configuration.DB.Where("invoice_id = ?", paymentRequest.InvoiceID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	if invoice.PaymentStatus == "Paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice already paid"})
		return
	}

	// Update payment status to "Paid" for offline payment
	invoice.PaymentStatus = "Paid"
	invoice.PaymentMethod = "Offline"
	if err := configuration.DB.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	// Update corresponding appointment status
	var appointment models.Appointment
	if err := configuration.DB.Where("appointment_id = ?", invoice.AppointmentID).First(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointment"})
		return
	}

	// Update appointment status to "confirmed" and payment status to "paid"
	appointment.BookingStatus = "confirmed"
	appointment.PaymentStatus = "paid"
	if err := configuration.DB.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Invoice payment successful",
		"invoice": invoice,
	})
}

type PageVariable struct {
	AppointmentID string
}

func MakePaymentOnline(c *gin.Context) {

	invoiceID := c.Query("id")
	id, err := strconv.Atoi(invoiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
	}

	var invoice models.Invoice
	if err := configuration.DB.First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "Failed",
				"message": "Invoice Not Found",
				"data":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch the invoice",
		})
		return
	}

	// Check if the invoice is already paid
	if invoice.PaymentStatus == "Paid" {
		c.JSON(400, gin.H{"error": "Invoice is already paid"})
		return
	}

	razorpayment := &models.RazorPay{
		InvoiceID:  uint(invoice.InvoiceID),
		AmountPaid: invoice.TotalAmount,
	}
	
	razorpayment.RazorPaymentID = generateUniqueID()
	if err := configuration.DB.Create(&razorpayment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create razor payment"})
		return
	}

	amountInPaisa := invoice.TotalAmount * 100
	razorpayClient := razorpay.NewClient(os.Getenv("RazorPay_key_id"), os.Getenv("RazorPay_key_secret"))

	data := map[string]interface{}{
		"amount":   amountInPaisa,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	body, err := razorpayClient.Order.Create(data, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Faied to create razorpay orer"})
	}

	value := body["id"]
	str := value.(string)

	homepagevariables := PageVariable{
		AppointmentID: str,
	}

	c.HTML(http.StatusOK, "payment.html", gin.H{
		"invoiceID":     id,
		"totalPrice":    amountInPaisa / 100,
		"total":         amountInPaisa,
		"appointmentID": homepagevariables.AppointmentID,
	})
}

func generateUniqueID() string {
	// Generate a Version 4 (random) UUID
	id := uuid.New()
	return id.String()
}

func SuccessPage(c *gin.Context) {
	paymentID := c.Query("bookID")
	fmt.Println(paymentID)
	var invoice models.Invoice
	if err := configuration.DB.First(&invoice, paymentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to fetch the invoice",
		})
		return
	}
	fmt.Printf("%+v\n", invoice)

	if invoice.PaymentStatus == "Pending" {
		if err := configuration.DB.Model(&invoice).Update("payment_status", "Paid").Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to update the invoice paymet status",
			})
			return
		}
	}

	razorPayment := models.RazorPay{
		InvoiceID:      uint(invoice.InvoiceID),
		RazorPaymentID: generateUniqueID(),
		AmountPaid:     invoice.TotalAmount,
	}

	if err := configuration.DB.Create(&razorPayment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to create RazorPay Payment",
		})
	}

	// Update appointment status in appointment table
	if invoice.AppointmentID != 0 {
		if err := configuration.DB.Model(&models.Appointment{}).Where("appointment_id = ?", invoice.AppointmentID).Updates(map[string]interface{}{"booking_status": "confirmed", "payment_status": "paid"}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to update the appointment status",
			})
			return
		}
	}
	
	c.HTML(http.StatusOK, "success.html", gin.H{
		"paymentID":   razorPayment.RazorPaymentID,
		"amountPaid": invoice.TotalAmount,
		"invoiceID":   invoice.InvoiceID,
	})
}

