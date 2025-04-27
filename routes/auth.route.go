package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ofojichigozie/hms-go-backend/controllers"
	"github.com/ofojichigozie/hms-go-backend/middleware"
	"github.com/ofojichigozie/hms-go-backend/repositories"
	"github.com/ofojichigozie/hms-go-backend/services"
	"gorm.io/gorm"
)

func AuthRoute(r *gin.Engine, DB *gorm.DB) {
	staffRepository := repositories.NewStaffRepository(DB)
	staffService := services.NewStaffService(staffRepository)
	authController := controllers.NewAuthController(staffService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.POST("/refresh", authController.RefreshToken)
		}
	}
}
