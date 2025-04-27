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

func ClinicalNoteRoutes(r *gin.Engine, DB *gorm.DB) {
	clinicalNoteRepository := repositories.NewClinicalNoteRepository(DB)
	appointmentRepository := repositories.NewAppointmentRepository(DB)
	patientRepository := repositories.NewPatientRepository(DB)
	noteService := services.NewClinicalNoteService(clinicalNoteRepository,
		appointmentRepository, patientRepository)
	noteController := controllers.NewClinicalNoteController(noteService)

	roles := constants.Roles

	noteGroup := r.Group("/clinical-notes")
	noteGroup.Use(middleware.AuthMiddleware())
	{
		doctorRoutes := noteGroup.Group("")
		doctorRoutes.Use(middleware.RoleMiddleware([]string{roles.DOCTOR}))
		{
			doctorRoutes.POST("", noteController.CreateNote)
			doctorRoutes.PATCH("/:id", noteController.UpdateNote)
			doctorRoutes.DELETE("/:id", noteController.DeleteNote)
		}

		staffRoutes := noteGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware([]string{roles.DOCTOR, roles.RECEPTIONIST}))
		{
			staffRoutes.GET("/:id", noteController.GetNoteByID)
			staffRoutes.GET("/patient/:patientId", noteController.GetNotesByPatientID)
		}
	}
}
