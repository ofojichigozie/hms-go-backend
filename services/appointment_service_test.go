package services

import (
	"errors"
	"testing"
	"time"

	"github.com/ofojichigozie/hms-go-backend/constants"
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateAppointment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		input := models.CreateAppointmentInput{
			PatientID:  1,
			Department: "cardiology",
			Duration:   30,
			Reason:     "Regular checkup",
		}

		createdBy := uint(2) // receptionist ID

		expectedPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(expectedPatient, nil)
		mockAppointmentRepo.On("Create", mock.AnythingOfType("*models.Appointment")).Return(nil).Run(func(args mock.Arguments) {
			appointment := args.Get(0).(*models.Appointment)
			assert.Equal(t, input.PatientID, appointment.PatientID)
			assert.Equal(t, createdBy, appointment.ReceptionistID)
			assert.Equal(t, input.Department, appointment.Department)
			assert.Equal(t, input.Duration, appointment.Duration)
			assert.Equal(t, input.Reason, appointment.Reason)
			assert.Equal(t, constants.AppointmentStatus.SCHEDULED, appointment.Status)
			assert.WithinDuration(t, time.Now(), appointment.ScheduledAt, time.Second)
		})

		result, err := service.CreateAppointment(input, createdBy)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockPatientRepo.AssertExpectations(t)
		mockAppointmentRepo.AssertExpectations(t)
	})

	t.Run("PatientNotFound", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		input := models.CreateAppointmentInput{
			PatientID:  1,
			Department: "cardiology",
		}

		createdBy := uint(2)

		mockPatientRepo.On("FindByID", uint(1)).Return(&models.Patient{}, errors.New("not found"))

		result, err := service.CreateAppointment(input, createdBy)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "associated patient record not found", err.Error())
		mockPatientRepo.AssertExpectations(t)
		mockAppointmentRepo.AssertNotCalled(t, "Create")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		input := models.CreateAppointmentInput{
			PatientID:  1,
			Department: "cardiology",
		}

		createdBy := uint(2)

		expectedPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(expectedPatient, nil)
		mockAppointmentRepo.On("Create", mock.AnythingOfType("*models.Appointment")).Return(errors.New("database error"))

		result, err := service.CreateAppointment(input, createdBy)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockPatientRepo.AssertExpectations(t)
		mockAppointmentRepo.AssertExpectations(t)
	})
}

func TestGetAllAppointments(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		filters := map[string]interface{}{"status": "scheduled"}
		expectedAppointments := []models.Appointment{
			{
				PatientID:  1,
				Department: "cardiology",
				Status:     "scheduled",
			},
		}

		mockAppointmentRepo.On("FindAll", filters).Return(expectedAppointments, nil)

		result, err := service.GetAllAppointments(filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedAppointments, result)
		mockAppointmentRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		filters := map[string]interface{}{}
		mockAppointmentRepo.On("FindAll", filters).Return([]models.Appointment{}, errors.New("database error"))

		result, err := service.GetAllAppointments(filters)

		assert.Error(t, err)
		assert.Equal(t, []models.Appointment{}, result)
		assert.Equal(t, "database error", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
	})
}

func TestGetAppointmentByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		expectedAppointment := &models.Appointment{
			Model:      gorm.Model{ID: 1},
			PatientID:  1,
			Department: "cardiology",
		}

		mockAppointmentRepo.On("FindByID", uint(1)).Return(expectedAppointment, nil)

		result, err := service.GetAppointmentByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedAppointment, result)
		mockAppointmentRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		mockAppointmentRepo.On("FindByID", uint(1)).Return(&models.Appointment{}, errors.New("not found"))

		result, err := service.GetAppointmentByID(1)

		assert.Error(t, err)
		assert.Equal(t, &models.Appointment{}, result)
		mockAppointmentRepo.AssertExpectations(t)
	})
}

func TestUpdateAppointment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		existingAppointment := &models.Appointment{
			Model:      gorm.Model{ID: 1},
			PatientID:  1,
			Department: "cardiology",
			Status:     "scheduled",
		}

		newDoctorID := uint(3)
		newStatus := "completed"
		newReason := "Patient recovered well"

		input := models.UpdateAppointmentInput{
			DoctorID: &newDoctorID,
			Status:   &newStatus,
			Reason:   &newReason,
		}

		updatedBy := uint(2)

		mockAppointmentRepo.On("FindByID", uint(1)).Return(existingAppointment, nil)
		mockAppointmentRepo.On("Update", mock.AnythingOfType("*models.Appointment")).Return(nil)

		result, err := service.UpdateAppointment(1, input, updatedBy)

		assert.NoError(t, err)
		assert.Equal(t, newDoctorID, *result.DoctorID)
		assert.Equal(t, newStatus, result.Status)
		assert.Equal(t, newReason, result.Reason)
		assert.Equal(t, updatedBy, result.UpdatedBy)
		mockAppointmentRepo.AssertExpectations(t)
	})

	t.Run("AppointmentNotFound", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		mockAppointmentRepo.On("FindByID", uint(1)).Return(&models.Appointment{}, errors.New("not found"))

		input := models.UpdateAppointmentInput{
			Status: stringPtr("completed"),
		}

		result, err := service.UpdateAppointment(1, input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "appointment not found", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockAppointmentRepo.AssertNotCalled(t, "Update")
	})

	t.Run("RepositoryUpdateError", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		existingAppointment := &models.Appointment{
			Model:      gorm.Model{ID: 1},
			PatientID:  1,
			Department: "cardiology",
		}

		newStatus := "completed"
		input := models.UpdateAppointmentInput{
			Status: &newStatus,
		}

		mockAppointmentRepo.On("FindByID", uint(1)).Return(existingAppointment, nil)
		mockAppointmentRepo.On("Update", mock.AnythingOfType("*models.Appointment")).Return(errors.New("update failed"))

		result, err := service.UpdateAppointment(1, input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "update failed", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
	})
}

func TestDeleteAppointment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		existingAppointment := &models.Appointment{
			Model:     gorm.Model{ID: 1},
			PatientID: 1,
			Status:    "scheduled",
		}

		mockAppointmentRepo.On("FindByID", uint(1)).Return(existingAppointment, nil)
		mockAppointmentRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeleteAppointment(1)

		assert.NoError(t, err)
		mockAppointmentRepo.AssertExpectations(t)
	})

	t.Run("AppointmentNotFound", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		mockAppointmentRepo.On("FindByID", uint(1)).Return(&models.Appointment{}, errors.New("appointment not found"))

		err := service.DeleteAppointment(1)

		assert.Error(t, err)
		assert.Equal(t, "appointment not found", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockAppointmentRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("CompletedAppointment", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		existingAppointment := &models.Appointment{
			Model:     gorm.Model{ID: 1},
			PatientID: 1,
			Status:    constants.AppointmentStatus.COMPLETED,
		}

		mockAppointmentRepo.On("FindByID", uint(1)).Return(existingAppointment, nil)

		err := service.DeleteAppointment(1)

		assert.Error(t, err)
		assert.Equal(t, "can't delete a completed appointment having clinical note(s)", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockAppointmentRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("RepositoryDeleteError", func(t *testing.T) {
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewAppointmentService(mockAppointmentRepo, mockPatientRepo)

		existingAppointment := &models.Appointment{
			Model:     gorm.Model{ID: 1},
			PatientID: 1,
			Status:    "scheduled",
		}

		mockAppointmentRepo.On("FindByID", uint(1)).Return(existingAppointment, nil)
		mockAppointmentRepo.On("Delete", uint(1)).Return(errors.New("delete failed"))

		err := service.DeleteAppointment(1)

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
	})
}

func stringPtr(s string) *string {
	return &s
}
