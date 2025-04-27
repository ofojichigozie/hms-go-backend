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

func AppointmentRoutes(r *gin.Engine, DB *gorm.DB) {
	staffRepository := repositories.NewStaffRepository(DB)
	appointmentRepository := repositories.NewAppointmentRepository(DB)
	appointmentService := services.NewAppointmentService(appointmentRepository, staffRepository)
	appointmentController := controllers.NewAppointmentController(appointmentService)

	roles := constants.Roles

	appointmentGroup := r.Group("/appointments")
	appointmentGroup.Use(middleware.AuthMiddleware())
	{
		receptionistRoutes := appointmentGroup.Group("")
		receptionistRoutes.Use(middleware.RoleMiddleware([]string{roles.RECEPTIONIST}))
		{
			receptionistRoutes.POST("", appointmentController.CreateAppointment)
			receptionistRoutes.PATCH("/:id", appointmentController.UpdateAppointment)
			receptionistRoutes.DELETE("/:id", appointmentController.DeleteAppointment)
		}

		doctorRoutes := appointmentGroup.Group("")
		doctorRoutes.Use(middleware.RoleMiddleware([]string{roles.DOCTOR}))
		{
			doctorRoutes.PATCH("/:id", appointmentController.UpdateAppointment)
		}

		staffRoutes := appointmentGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware([]string{roles.RECEPTIONIST, roles.DOCTOR}))
		{
			staffRoutes.GET("", appointmentController.GetAllAppointments)
			staffRoutes.GET("/:id", appointmentController.GetAppointmentByID)
		}
	}
}
