package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ofojichigozie/hms-go-backend/responses"
	"github.com/ofojichigozie/hms-go-backend/services"
	"github.com/ofojichigozie/hms-go-backend/utils"
)

type AuthController struct {
	staffService services.StaffService
}

func NewAuthController(staffService services.StaffService) *AuthController {
	return &AuthController{staffService}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var body LoginInput

	if err := ctx.ShouldBindJSON(&body); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	email := strings.ToLower(body.Email)

	staff, err := ac.staffService.GetStaffByEmail(email)
	if err != nil || staff == nil {
		responses.Error(ctx, http.StatusUnauthorized,
			"Authentication failed", "Invalid email or password")
		return
	}

	if err := utils.VerifyPassword(staff.PasswordHash, body.Password); err != nil {
		responses.Error(ctx, http.StatusUnauthorized,
			"Authentication failed", "Invalid email or password")
		return
	}

	if !staff.IsActive {
		responses.Error(ctx, http.StatusForbidden,
			"Account inactive", "Please contact administrator")
		return
	}

	tokens, err := utils.GenerateTokenPair(staff.ID, staff.Role)
	if err != nil {
		responses.Error(ctx, http.StatusInternalServerError, "Token generation failed", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Login successful", gin.H{
		"accessToken":  tokens["accessToken"],
		"refreshToken": tokens["refreshToken"],
		"staff": gin.H{
			"id":         staff.ID,
			"employeeId": staff.EmployeeID,
			"firstName":  staff.FirstName,
			"lastName":   staff.LastName,
			"email":      staff.Email,
			"role":       staff.Role,
			"isActive":   staff.IsActive,
		},
	})
}

func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		responses.Error(ctx, http.StatusUnauthorized, "Authorization header required", nil)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		responses.Error(ctx, http.StatusUnauthorized, "Invalid authorization format", nil)
		return
	}

	newTokens, err := utils.RefreshToken(tokenParts[1])
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Token refresh failed", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Token refreshed successfully", gin.H{
		"accessToken":  newTokens["accessToken"],
		"refreshToken": newTokens["refreshToken"],
	})
}
