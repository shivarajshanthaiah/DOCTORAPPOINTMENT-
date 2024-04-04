package routes

import (
	"doctorAppointment/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes() *gin.Engine {
	//creates a new Gin engine instance with default configurations
	r := gin.Default()

	//user routers
	r.POST("/users/login", controllers.PatientLogin)
	r.POST("/users/signup", controllers.PatientSignup)
	r.POST("/users/verify", controllers.UserOtpVerify)
	r.GET("/users/logout", controllers.PatientLogout)

	r.POST("/admin/login", controllers.AdminLogin)
	r.POST("/admin/logout", controllers.AdminLogout)

	//doctor routes

	return r
}
