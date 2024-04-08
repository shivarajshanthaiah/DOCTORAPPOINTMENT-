package routes

import (
	"doctorAppointment/authentication"
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


	//Admin routes
	
	r.POST("/admin/login", controllers.AdminLogin)

	admin := r.Group("/admin")
	admin.Use(authentication.AdminAuthMiddleware())
	{
		admin.POST("/logout", controllers.AdminLogout)
		admin.GET("/view/hospitals", controllers.ViewHospitals)
		admin.POST("/add/hospital", controllers.AddHospital)
		admin.GET("/search/hospital/:id", controllers.SearchHospital)
		admin.PATCH("/update/hospital/:id", controllers.UpdateHospital)
		admin.POST("/remove/hospital/:id", controllers.RemoveHospital)
		

	}

	return r
}
