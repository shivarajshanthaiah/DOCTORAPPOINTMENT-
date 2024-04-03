package configuration

import (
	"doctorAppointment/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// hold connectioin to db
var DB *gorm.DB

// initializing db connection
func ConfigDB() {

	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DB")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	DB.AutoMigrate(&models.Appointment{}, 
		&models.AppointmentHistory{}, 
		&models.Doctor{}, 
		&models.Hospital{}, 
		&models.Patient{}, 
		&models.Payment{}, 
		&models.Prescription{}, 
		&models.Slot{},
	)

}
