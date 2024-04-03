package main

import (
	"doctorAppointment/configuration"
	"doctorAppointment/routes"
)
func Init(){
	configuration.ConfigDB()
	configuration.InitRedis()
	// configuration.Loadenv()
}

	func main() {
		//Perform application initialization
		Init()
		r := routes.UserRoutes()
	
		//Run the engine the port 3000
		if err := r.Run(":3000"); err != nil {
			panic(err)
		}
	
	}
