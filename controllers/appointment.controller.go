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

type AppointmentController struct {
	appointmentService services.AppointmentService
}

func NewAppointmentController(appointmentService services.AppointmentService) *AppointmentController {
	return &AppointmentController{appointmentService}
}

func (ac *AppointmentController) CreateAppointment(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	var input models.CreateAppointmentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	appointment, err := ac.appointmentService.CreateAppointment(input, currentStaff.ID)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to create appointment", err.Error())
		return
	}

	responses.Success(ctx, http.StatusCreated, "Appointment created successfully", appointment)
}

func (ac *AppointmentController) GetAllAppointments(ctx *gin.Context) {
	filters := make(map[string]interface{})

	if patientID := ctx.Query("patientId"); patientID != "" {
		filters["patient_id"] = patientID
	}
	if doctorID := ctx.Query("doctorId"); doctorID != "" {
		filters["doctor_id"] = doctorID
	}
	if department := ctx.Query("department"); department != "" {
		filters["department"] = department
	}
	if status := ctx.Query("status"); status != "" {
		filters["status"] = status
	}

	appointments, err := ac.appointmentService.GetAllAppointments(filters)
	if err != nil {
		responses.Error(ctx, http.StatusInternalServerError, "Failed to fetch appointments", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Appointments retrieved successfully", appointments)
}

func (ac *AppointmentController) GetAppointmentByID(ctx *gin.Context) {
	appointmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid appointment ID", "Appointment ID must be a positive integer")
		return
	}

	appointment, err := ac.appointmentService.GetAppointmentByID(uint(appointmentID))
	if err != nil {
		responses.Error(ctx, http.StatusNotFound, "Appointment not found", nil)
		return
	}

	responses.Success(ctx, http.StatusOK, "Appointment retrieved successfully", appointment)
}

func (ac *AppointmentController) UpdateAppointment(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	appointmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid appointment ID", "Appointment ID must be a positive integer")
		return
	}

	var input models.UpdateAppointmentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	appointment, err := ac.appointmentService.UpdateAppointment(uint(appointmentID), input, currentStaff.ID)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to update appointment", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Appointment updated successfully", appointment)
}

func (ac *AppointmentController) DeleteAppointment(ctx *gin.Context) {
	appointmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid appointment ID", "Appointment ID must be a positive integer")
		return
	}

	err = ac.appointmentService.DeleteAppointment(uint(appointmentID))
	if err != nil {
		responses.Error(ctx, http.StatusNotFound, "Appointment not found", nil)
		return
	}

	responses.Success(ctx, http.StatusOK, "Appointment deleted successfully", nil)
}
