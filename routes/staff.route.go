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

func StaffRoutes(r *gin.Engine, DB *gorm.DB) {
	staffRepository := repositories.NewStaffRepository(DB)
	staffService := services.NewStaffService(staffRepository)
	staffController := controllers.NewStaffController(staffService)

	roles := constants.Roles

	staffGroup := r.Group("/staff")
	staffGroup.Use(middleware.AuthMiddleware())
	{
		adminRoutes := staffGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware([]string{roles.ADMIN}))
		{
			adminRoutes.POST("", staffController.CreateStaff)
			staffGroup.GET("", staffController.GetAllStaff)
			adminRoutes.PATCH("/:id", staffController.UpdateStaff)
			adminRoutes.DELETE("/:id", staffController.DeleteStaff)
		}

		staffGroup.GET("/:id", staffController.GetStaffByID)

	}
}
