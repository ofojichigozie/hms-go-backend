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

type ClinicalNoteController struct {
	clinicalNoteService services.ClinicalNoteService
}

func NewClinicalNoteController(clinicalNoteService services.ClinicalNoteService) *ClinicalNoteController {
	return &ClinicalNoteController{clinicalNoteService}
}

func (c *ClinicalNoteController) CreateNote(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	var input models.CreateNoteInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	note, err := c.clinicalNoteService.CreateNote(input, currentStaff.ID)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to create clinical note", err.Error())
		return
	}

	responses.Success(ctx, http.StatusCreated, "Clinical note created successfully", note)
}

func (c *ClinicalNoteController) GetNoteByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid note ID", "Note ID must be a positive integer")
		return
	}

	note, err := c.clinicalNoteService.GetNoteByID(uint(id))
	if err != nil {
		responses.Error(ctx, http.StatusNotFound, "Clinical note not found", nil)
		return
	}

	responses.Success(ctx, http.StatusOK, "Clinical note retrieved successfully", note)
}

func (c *ClinicalNoteController) GetNotesByPatientID(ctx *gin.Context) {
	patientID, err := strconv.ParseUint(ctx.Param("patientId"), 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid patient ID", "Patient ID must be a positive integer")
		return
	}

	notes, err := c.clinicalNoteService.GetNotesByPatientID(uint(patientID))
	if err != nil {
		responses.Error(ctx, http.StatusInternalServerError, "Failed to retrieve notes", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Clinical notes retrieved successfully", notes)
}

func (c *ClinicalNoteController) UpdateNote(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	clinicalNoteId, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid note ID", "Note ID must be a positive integer")
		return
	}

	var input models.UpdateNoteInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	note, err := c.clinicalNoteService.UpdateNote(uint(clinicalNoteId), input, currentStaff.ID)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to update clinical note", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Clinical note updated successfully", note)
}

func (c *ClinicalNoteController) DeleteNote(ctx *gin.Context) {
	currentStaff, err := middleware.GetCurrentStaff(ctx)
	if err != nil {
		responses.Error(ctx, http.StatusUnauthorized, "Authentication required", err.Error())
		return
	}

	clinicalNoteId, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest,
			"Invalid note ID", "Note ID must be a positive integer")
		return
	}

	err = c.clinicalNoteService.DeleteNote(uint(clinicalNoteId), currentStaff.ID)
	if err != nil {
		responses.Error(ctx, http.StatusBadRequest, "Failed to delete clinical note", err.Error())
		return
	}

	responses.Success(ctx, http.StatusOK, "Clinical note deleted successfully", nil)
}
