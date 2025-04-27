package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ofojichigozie/hms-go-backend/constants"
	"github.com/ofojichigozie/hms-go-backend/controllers"
	"github.com/ofojichigozie/hms-go-backend/middleware"
	"github.com/ofojichigozie/hms-go-backend/repositories"
	"github.com/ofojichigozie/hms-go-backend/services"
	"gorm.io/gorm"
)

func PatientRoutes(r *gin.Engine, DB *gorm.DB) {
	patientRepository := repositories.NewPatientRepository(DB)
	staffRepository := repositories.NewStaffRepository(DB)
	patientService := services.NewPatientService(patientRepository, staffRepository)
	patientController := controllers.NewPatientController(patientService)

	roles := constants.Roles

	patientGroup := r.Group("/patients")
	patientGroup.Use(middleware.AuthMiddleware())
	{
		receptionistRoutes := patientGroup.Group("")
		receptionistRoutes.Use(middleware.RoleMiddleware([]string{roles.RECEPTIONIST}))
		{
			receptionistRoutes.POST("", patientController.CreatePatient)
			receptionistRoutes.DELETE("/:id", patientController.DeletePatient)
			receptionistRoutes.PATCH("/:id", patientController.UpdatePatient)
		}

		staffRoutes := patientGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware([]string{roles.RECEPTIONIST, roles.DOCTOR}))
		{
			staffRoutes.GET("", patientController.GetAllPatients)
			staffRoutes.GET("/:id", patientController.GetPatientByID)
		}
	}
}
