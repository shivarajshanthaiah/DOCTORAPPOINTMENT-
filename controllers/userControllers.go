package controllers

import (
	"doctorAppointment/authentication"
	"doctorAppointment/configuration"
	"doctorAppointment/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User Login
func PatientLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Doctor Appointment Booking. Please login...!"})
}

func PatientSignup(c *gin.Context) {
	var patient models.Patient
	if err := c.BindJSON(&patient); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(patient)

	if err := configuration.DB.Where("phone = ?", patient.Phone).First(&models.Patient{}).Error; err == nil {
		//Genereate Token
		token, err := authentication.GeneratePatienttoken(patient.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Token Generated Successfully", "token": token})
		return

	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error exist"})
		return
	}

	err := SendOTP(patient.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP", "data": err.Error()})
		return
	}

	key := fmt.Sprintf("user:%s", patient.Phone)
	err = configuration.SetRedis(key, patient.Phone, time.Minute*5)
	if err != nil{
		fmt.Println("Error setting user in Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Status": false, "Data": nil, "Message":"Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message":"Otp generated successfully. Proceed to verification page>>>"})

}

func SendOTP(phoneNumber string) error {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTHTOKEN")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	//create SMS message
	from := os.Getenv("TWILIO_PHONENUMBER")
	params := verify.CreateVerificationParams{}
	params.SetTo("+918762334325")
	params.SetChannel("sms")
	println(from)
	response, err := client.VerifyV2.CreateVerification(os.Getenv("TWILIO_SERVIES_ID"), &params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(response)
	return nil
}

func OTPverify(c *gin.Context) {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTHTOKEN")

	var OTPverify models.VerifyOTP
	if err := c.BindJSON(&OTPverify); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": false, "Data": nil, "Message": err.Error()})
		return
	}


	if OTPverify.Otp == "" {
		c.JSON(http.StatusOK, gin.H{"Status": true, "Message": "OTP verified successfully"})
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	params := verify.CreateVerificationCheckParams{}
	params.SetTo("+918762334325")
	params.SetCode(OTPverify.Otp)

	//send twilio verification check
	response, err := client.VerifyV2.CreateVerificationCheck(os.Getenv("TWILIO_SERVIES_ID"), &params)
	if err != nil {
		fmt.Println("err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Status": false, "Data": nil, "Message": "error in veifying provided OTP"})
		return
	} else if *response.Status != "approved" {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": false, "Data": nil, "Message": "Wrong OTP provoded"})
		return
	}

	key := fmt.Sprintf("user:%s", OTPverify.Phone)
	value, err := configuration.GetRedis(key)
	if err != nil {
		fmt.Println("Error checking user in Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "Data": nil, "Message": "Internal server error"})
		return
	}

	var patient models.Patient
	patient.Phone = value
	patient.Name = value
	patient.Age = value
	patient.Address = value
	patient.Email = value
	patient.Gender= value

	err = configuration.DB.Create(&patient).Error
	if err != nil {
		fmt.Println("Error creating Patient:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Status": false, "Data": nil, "Message": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": true, "Message": "OTP verified successfully"})
}

// User logout
func PatientLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "You are Sucessfully logged out"})
}
