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

func GetInvoice(c *gin.Context) {
	var invoice []models.Invoice
	if err := configuration.DB.Find(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occured while receiving the invoice",
		})
		return
	}
	c.JSON(http.StatusOK, invoice)
}

// To make payment offline
func PayInvoiceOffline(c *gin.Context) {
	// Struct to hold the payment request parameters
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

// PageVariable struct holds data to be passed to the HTML template.
type PageVariable struct {
	AppointmentID string
}

// Function for processing online payments
func MakePaymentOnline(c *gin.Context) {

	invoiceID := c.Query("id")
	// Convert the invoice ID from string to integer
	id, err := strconv.Atoi(invoiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
	}

	// Retrieve the invoice corresponding to the provided ID from the database
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

	// Create a RazorPay payment record in the database with the invoice ID and total amount.
	razorpayment := &models.RazorPay{
		InvoiceID:  uint(invoice.InvoiceID),
		AmountPaid: invoice.TotalAmount,
	}
	
	razorpayment.RazorPaymentID = generateUniqueID()
	if err := configuration.DB.Create(&razorpayment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create razor payment"})
		return
	}

	// Convert total amount to paisa (multiply by 100) for RazorPay API
	amountInPaisa := invoice.TotalAmount * 100
	razorpayClient := razorpay.NewClient(os.Getenv("RazorPay_key_id"), os.Getenv("RazorPay_key_secret"))

	// Prepare data for creating a RazorPay order.
	data := map[string]interface{}{
		"amount":   amountInPaisa,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	// Create a RazorPay order using the RazorPay API
	body, err := razorpayClient.Order.Create(data, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Faied to create razorpay orer"})
	}

	// Extract the order ID from the response body returned by the RazorPay API
	value := body["id"]
	str := value.(string)

	// Create an instance of the PageVariable struct to hold data for the HTML template
	homepagevariables := PageVariable{
		AppointmentID: str,
	}

	// Render the payment.html template, passing invoice ID, total price, total amount, and appointment ID as template variables.
	c.HTML(http.StatusOK, "payment.html", gin.H{
		"invoiceID":     id,
		"totalPrice":    amountInPaisa / 100,
		"total":         amountInPaisa,
		"appointmentID": homepagevariables.AppointmentID,
	})
}

// generateUniqueID generates a unique ID using UUID (Universally Unique Identifier).
func generateUniqueID() string {
	// Generate a Version 4 (random) UUID
	id := uuid.New()
	return id.String()
}


//Function to display success page after successfull payment
func SuccessPage(c *gin.Context) {
	paymentID := c.Query("bookID")
	fmt.Println(paymentID)

	// Fetch the invoice corresponding to the provided payment ID from the database
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

	// Update payment method to "online"
    if err := configuration.DB.Model(&invoice).Update("payment_method", "online").Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "Error": "Failed to update the payment method",
        })
        return
    }

	// Create a record of the RazorPay payment in the database
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
	
	// Render the success page template, passing payment ID, amount paid, and invoice ID as template variables
	c.HTML(http.StatusOK, "success.html", gin.H{
		"paymentID":   razorPayment.RazorPaymentID,
		"amountPaid": invoice.TotalAmount,
		"invoiceID":   invoice.InvoiceID,
	})
}

