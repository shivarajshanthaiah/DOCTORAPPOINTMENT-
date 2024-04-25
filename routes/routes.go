package routes

import (
	"doctorAppointment/authentication"
	"doctorAppointment/controllers"
	adminControllers "doctorAppointment/controllers/adminControllers"
	"doctorAppointment/controllers/doctorControllers"
	"doctorAppointment/controllers/userControllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes() *gin.Engine {
	//creates a new Gin engine instance with default configurations
	r := gin.Default()

	//user routers
	r.POST("/users/login", userControllers.PatientLogin)
	r.POST("/users/signup", userControllers.PatientSignup)
	r.POST("/users/verify", userControllers.UserOtpVerify)
	r.GET("/pay/invoice/online", controllers.MakePaymentOnline)
	r.GET("/payment/success", controllers.SuccessPage)

	user := r.Group("/user")
	user.Use(authentication.PatientAuthMiddleware())
	{
		user.GET("/doctors/:doctor_id/available-slots", userControllers.GetAvailableTimeSlots)
		user.GET("/logout", userControllers.PatientLogout)
		user.GET("/doctor/:specialization", userControllers.GetDoctorsBySpeciality)
		user.POST("/book/appointment", userControllers.BookAppointment)
		user.POST("/pay/invoice/offline", controllers.PayInvoiceOffline)
		user.GET("/wallet/:userid", controllers.Wallet)
		user.POST("/cancel/appointment/:id", controllers.CancelAppointment)
		user.GET("/appointment/history/:id", userControllers.GetAppointmenentHistory)

	}

	//Admin routes

	r.POST("/admin/login", adminControllers.AdminLogin)

	admin := r.Group("/admin")
	admin.Use(authentication.AdminAuthMiddleware())
	{
		admin.POST("/logout", adminControllers.AdminLogout)
		admin.GET("/view/hospitals", adminControllers.ViewHospitals)
		admin.POST("/add/hospital", adminControllers.AddHospital)
		admin.GET("/search/hospital/:id", adminControllers.SearchHospital)
		admin.PATCH("/update/hospital/:id", adminControllers.UpdateHospital)
		admin.POST("/remove/hospital/:id", adminControllers.RemoveHospital)
		admin.GET("/view/deleted/hospitals", adminControllers.ViewDeletedHospitals)
		admin.GET("/view/Active/hospitals", adminControllers.ViewActiveHospitals)
		admin.POST("/verify/doctor/:id", adminControllers.UpdateDoctor)
		admin.GET("/view/verified/doctors", adminControllers.ViewVerifiedDoctors)
		admin.GET("/view/doctor/:id", adminControllers.GetDoctorByID)
		admin.GET("/view/doctors/:specialization", adminControllers.GetDoctorBySpeciality)
		admin.GET("/view/notVerified/doctors", adminControllers.ViewNotVerifiedDoctors)
		admin.GET("/view/verified/approved/doctors", adminControllers.ViewVerifiedApprovedDoctors)
		admin.GET("/view/verified/notApproved/doctors", adminControllers.ViewVerifiedNotApprovedDoctors)
		admin.GET("/view/invoice", controllers.GetInvoice)
		admin.GET("/total/appointments", adminControllers.GetBookingStatusCounts)
		admin.GET("/doctor-wise/bookings", adminControllers.GetDoctorWiseBookings)
		admin.GET("/department-wise/bookings", adminControllers.GetDepartmentWiseBookings)
		admin.GET("/total/revenue", adminControllers.GetTotalRevenue)
		admin.GET("/revenue/startdate", adminControllers.GetSpecificRevenue)
	}

	//Doctor routes
	r.POST("doctor/signup", doctorControllers.Signup)
	r.POST("doctor/verify", doctorControllers.VerifyOTP)
	r.GET("view/hospitals", doctorControllers.ViewHospital)
	//r.POST("doctor/signup", doctorControllers.DoctorSignup)
	r.POST("/doctor/login", doctorControllers.DoctorLogin)

	doctors := r.Group("/doctor")
	doctors.Use(authentication.DoctorAuthMiddleware())
	{
		doctors.POST("/update/availability", doctorControllers.SaveAvailability)
		doctors.GET("/logout", doctorControllers.DoctorLogout)
		doctors.POST("/add/prescription", doctorControllers.AddPrescription)
		doctors.POST("/cancel/appointment/:id", controllers.CancelAppointment)
		doctors.GET("/appointment/history/:id", doctorControllers.GetAppHistory)
		doctors.GET("/appointment/:doctor_id/date", doctorControllers.GetDoctorAppointmentsByDate)
	}

	return r
}
