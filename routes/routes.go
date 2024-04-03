package routes

import (
	"doctorAppointment/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes() *gin.Engine {
	//creates a new Gin engine instance with default configurations
	r := gin.Default()

	//user routers
	r.GET("/users/login", controllers.PatientLogin)
	r.POST("/users/signup", controllers.PatientSignup)
	r.POST("/users/verify", controllers.OTPverify)
	r.GET("/users/logout", controllers.PatientLogout)
	


	//doctor routes
	
	return r
}
