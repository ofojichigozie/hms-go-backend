package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ofojichigozie/hms-go-backend/initializers"
	"github.com/ofojichigozie/hms-go-backend/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	if _, err := initializers.CreateInitialAdmin(initializers.DB); err != nil {
		log.Printf("Warning: Admin creation failed: %v", err)
	} else {
		log.Println("Admin initialization completed")
	}

	r := gin.Default()
	routes.AuthRoute(r, initializers.DB)
	routes.StaffRoutes(r, initializers.DB)
	routes.PatientRoutes(r, initializers.DB)
	routes.AppointmentRoutes(r, initializers.DB)
	routes.ClinicalNoteRoutes(r, initializers.DB)
	r.Run()
}
