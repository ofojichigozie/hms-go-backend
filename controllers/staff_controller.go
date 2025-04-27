package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ofojichigozie/hms-go-backend/middleware"
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/responses"
	"github.com/ofojichigozie/hms-go-backend/services"
)

type StaffController struct {
	staffService services.StaffService
}

func NewStaffController(staffService services.StaffService) *StaffController {
	return &StaffController{staffService}
}

func (sc *StaffController) CreateStaff(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	if currentStaff.Role != "admin" {
		responses.Error(ctx, http.StatusForbidden,
			"Permission denied", "Only admin can create staff accounts")
		return
	}

	var input models.CreateStaffInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	staff, err := sc.staffService.CreateStaff(input)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to create staff", err.Error())
		return
	}

	responses.Success(ctx, http.StatusCreated, "Staff created successfully", staff)
}

func (sc *StaffController) GetAllStaff(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	if currentStaff.Role != "admin" {
		responses.Error(ctx, http.StatusForbidden,
			"Permission denied", "Only admin can view all staff records")
		return
	}

	staffList, err := sc.staffService.GetAllStaff()
	if err != nil {
		responses.Error(ctx, http.StatusInternalServerError, "Failed to fetch staff", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Staff list retrieved successfully", staffList)
}

func (sc *StaffController) GetStaffByID(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	id := ctx.Param("id")
	staffID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid staff ID", "Staff ID must be a positive integer")
		return
	}

	if currentStaff.Role != "admin" && currentStaff.ID != uint(staffID) {
		responses.Error(ctx, http.StatusForbidden,
			"Permission denied", "You can only view your own staff record")
		return
	}

	staff, err := sc.staffService.GetStaffByID(uint(staffID))
	if err != nil {
		responses.Error(ctx, http.StatusNotFound, "Staff not found", nil)
		return
	}

	responses.Success(ctx, http.StatusOK, "Staff retrieved successfully", staff)
}

func (sc *StaffController) UpdateStaff(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	if currentStaff.Role != "admin" {
		responses.Error(ctx, http.StatusForbidden,
			"Permission denied", "Only admin can update staff records")
		return
	}

	id := ctx.Param("id")
	staffID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid staff ID", "Staff ID must be a positive integer")
		return
	}

	var input models.UpdateStaffInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	staff, err := sc.staffService.UpdateStaff(uint(staffID), input)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to update staff", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Staff updated successfully", staff)
}

func (sc *StaffController) DeleteStaff(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	if currentStaff.Role != "admin" {
		responses.Error(ctx, http.StatusForbidden,
			"Permission denied", "Only admin can delete staff records")
		return
	}

	id := ctx.Param("id")
	staffID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid staff ID", "Staff ID must be a positive integer")
		return
	}

	err = sc.staffService.DeleteStaff(uint(staffID))
	if err != nil {
		responses.Error(ctx, http.StatusNotFound, "Staff not found", nil)
		return
	}

	responses.Success(ctx, http.StatusOK, "Staff deleted successfully", nil)
}
