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

type PatientController struct {
	patientService services.PatientService
}

func NewPatientController(patientService services.PatientService) *PatientController {
	return &PatientController{patientService}
}

func (pc *PatientController) CreatePatient(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	var input models.CreatePatientInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	patient, err := pc.patientService.CreatePatient(input, currentStaff.ID)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to create patient record", err.Error())
		return
	}

	responses.Success(ctx, http.StatusCreated, "Patient created successfully", patient)
}

func (pc *PatientController) GetAllPatients(ctx *gin.Context) {
	filters := make(map[string]interface{})
	if registrationNumber := ctx.Query("registrationNumber"); registrationNumber != "" {
		filters["registration_number"] = registrationNumber
	}
	// TODO: Add more filter parameters from query string

	patients, err := pc.patientService.GetAllPatients(filters)
	if err != nil {
		responses.Error(ctx, http.StatusInternalServerError, "Failed to fetch patients", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Patients retrieved successfully", patients)
}

func (pc *PatientController) GetPatientByID(ctx *gin.Context) {
	id := ctx.Param("id")
	patientID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid patient ID", "Patient ID must be a positive integer")
		return
	}

	patient, err := pc.patientService.GetPatientByID(uint(patientID))
	if err != nil {
		responses.Error(ctx, http.StatusNotFound, "Patient not found", nil)
		return
	}

	responses.Success(ctx, http.StatusOK, "Patient retrieved successfully", patient)
}

func (pc *PatientController) UpdatePatient(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	id := ctx.Param("id")
	patientID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid patient ID", "Patient ID must be a positive integer")
		return
	}

	var input models.UpdatePatientInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	patient, err := pc.patientService.UpdatePatient(uint(patientID), input, currentStaff.ID)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to update patient profile", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Patient profile updated successfully", patient)
}

func (pc *PatientController) DeletePatient(ctx *gin.Context) {
	id := ctx.Param("id")
	patientID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid patient ID", "Patient ID must be a positive integer")
		return
	}

	err = pc.patientService.DeletePatient(uint(patientID))
	if err != nil {
		responses.Error(ctx, http.StatusNotFound, "Patient not found", nil)
		return
	}

	responses.Success(ctx, http.StatusOK, "Patient deleted successfully", nil)
}
