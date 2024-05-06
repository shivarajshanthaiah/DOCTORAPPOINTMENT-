package main

import (
	"doctorAppointment/configuration"
	"doctorAppointment/routes"
)
func Init(){
	configuration.ConfigDB()
	configuration.InitRedis()
}

	func main() {
		//Perform application initialization
		Init()
		r := routes.UserRoutes()
		r.LoadHTMLGlob("templates/*")
	
		//Run the engine in default port
		if err := r.Run(); err != nil {
			panic(err)
		}
	
	}
