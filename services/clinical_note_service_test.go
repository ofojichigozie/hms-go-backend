package services

import (
	"errors"
	"testing"

	"github.com/ofojichigozie/hms-go-backend/constants"
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateNote(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		input := models.CreateNoteInput{
			AppointmentID:        1,
			PresentingComplaints: "Headache",
			PastMedicalHistory:   "None",
			ClinicalDiagnosis:    "Migraine",
			TreatmentPlan:        "Rest and medication",
			Recommendation:       "Follow up in 2 weeks",
		}

		doctorID := uint(2)

		expectedAppointment := &models.Appointment{
			Model:     gorm.Model{ID: 1},
			PatientID: 1,
			Status:    constants.AppointmentStatus.SCHEDULED,
		}

		mockAppointmentRepo.On("FindByID", uint(1)).Return(expectedAppointment, nil)
		mockAppointmentRepo.On("Update", mock.AnythingOfType("*models.Appointment")).Return(nil)
		mockNoteRepo.On("Create", mock.AnythingOfType("*models.ClinicalNote")).Return(nil).Run(func(args mock.Arguments) {
			note := args.Get(0).(*models.ClinicalNote)
			assert.Equal(t, input.AppointmentID, note.AppointmentID)
			assert.Equal(t, expectedAppointment.PatientID, note.PatientID)
			assert.Equal(t, doctorID, note.DoctorID)
			assert.Equal(t, input.PresentingComplaints, note.PresentingComplaints)
			assert.Equal(t, input.PastMedicalHistory, note.PastMedicalHistory)
			assert.Equal(t, input.ClinicalDiagnosis, note.ClinicalDiagnosis)
			assert.Equal(t, input.TreatmentPlan, note.TreatmentPlan)
			assert.Equal(t, input.Recommendation, note.Recommendation)
		})

		result, err := service.CreateNote(input, doctorID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockAppointmentRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("AppointmentNotFound", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		input := models.CreateNoteInput{
			AppointmentID: 1,
		}

		mockAppointmentRepo.On("FindByID", uint(1)).Return(&models.Appointment{}, errors.New("not found"))

		result, err := service.CreateNote(input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "associated appointment record not found", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockNoteRepo.AssertNotCalled(t, "Create")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		input := models.CreateNoteInput{
			AppointmentID:        1,
			PresentingComplaints: "Headache",
		}

		expectedAppointment := &models.Appointment{
			Model:     gorm.Model{ID: 1},
			PatientID: 1,
		}

		mockAppointmentRepo.On("FindByID", uint(1)).Return(expectedAppointment, nil)
		mockNoteRepo.On("Create", mock.AnythingOfType("*models.ClinicalNote")).Return(errors.New("database error"))

		result, err := service.CreateNote(input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockAppointmentRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})
}

func TestGetNoteByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		expectedNote := &models.ClinicalNote{
			Model:                gorm.Model{ID: 1},
			AppointmentID:        1,
			PresentingComplaints: "Headache",
		}

		mockNoteRepo.On("FindByID", uint(1)).Return(expectedNote, nil)

		result, err := service.GetNoteByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedNote, result)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		mockNoteRepo.On("FindByID", uint(1)).Return(&models.ClinicalNote{}, errors.New("not found"))

		result, err := service.GetNoteByID(1)

		assert.Error(t, err)
		assert.Equal(t, &models.ClinicalNote{}, result)
		mockNoteRepo.AssertExpectations(t)
	})
}

func TestGetNoteByAppointmentID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		expectedNote := &models.ClinicalNote{
			Model:         gorm.Model{ID: 1},
			AppointmentID: 1,
		}

		mockNoteRepo.On("FindByAppointmentID", uint(1)).Return(expectedNote, nil)

		result, err := service.GetNoteByAppointmentID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedNote, result)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		mockNoteRepo.On("FindByAppointmentID", uint(1)).Return(&models.ClinicalNote{}, errors.New("not found"))

		result, err := service.GetNoteByAppointmentID(1)

		assert.Error(t, err)
		assert.Equal(t, &models.ClinicalNote{}, result)
		mockNoteRepo.AssertExpectations(t)
	})
}

func TestGetNotesByPatientID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		expectedNotes := []models.ClinicalNote{
			{
				Model:     gorm.Model{ID: 1},
				PatientID: 1,
			},
		}

		expectedPatient := &models.Patient{
			Model: gorm.Model{ID: 1},
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(expectedPatient, nil)
		mockNoteRepo.On("FindByPatientID", uint(1)).Return(expectedNotes, nil)

		result, err := service.GetNotesByPatientID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedNotes, result)
		mockPatientRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("PatientNotFound", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		mockPatientRepo.On("FindByID", uint(1)).Return(&models.Patient{}, errors.New("not found"))

		result, err := service.GetNotesByPatientID(1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "patient record not found", err.Error())
		mockPatientRepo.AssertExpectations(t)
		mockNoteRepo.AssertNotCalled(t, "FindByPatientID")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		expectedPatient := &models.Patient{
			Model: gorm.Model{ID: 1},
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(expectedPatient, nil)
		mockNoteRepo.On("FindByPatientID", uint(1)).Return([]models.ClinicalNote{}, errors.New("database error"))

		result, err := service.GetNotesByPatientID(1)

		assert.Error(t, err)
		assert.Equal(t, []models.ClinicalNote{}, result)
		assert.Equal(t, "database error", err.Error())
		mockPatientRepo.AssertExpectations(t)
		mockNoteRepo.AssertExpectations(t)
	})
}

func TestUpdateNote(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		existingNote := &models.ClinicalNote{
			Model:                gorm.Model{ID: 1},
			DoctorID:             2,
			PresentingComplaints: "Old complaint",
		}

		newComplaint := "Updated complaint"
		input := models.UpdateNoteInput{
			PresentingComplaints: &newComplaint,
		}

		mockNoteRepo.On("FindByID", uint(1)).Return(existingNote, nil)
		mockNoteRepo.On("Update", mock.AnythingOfType("*models.ClinicalNote")).Return(nil)

		result, err := service.UpdateNote(1, input, 2)

		assert.NoError(t, err)
		assert.Equal(t, newComplaint, result.PresentingComplaints)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("NoteNotFound", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		mockNoteRepo.On("FindByID", uint(1)).Return(&models.ClinicalNote{}, errors.New("not found"))

		input := models.UpdateNoteInput{
			PresentingComplaints: stringPtr("New complaint"),
		}

		result, err := service.UpdateNote(1, input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "clinical note not found", err.Error())
		mockNoteRepo.AssertExpectations(t)
		mockNoteRepo.AssertNotCalled(t, "Update")
	})

	t.Run("UnauthorizedUpdate", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		existingNote := &models.ClinicalNote{
			Model:    gorm.Model{ID: 1},
			DoctorID: 2,
		}

		mockNoteRepo.On("FindByID", uint(1)).Return(existingNote, nil)

		input := models.UpdateNoteInput{
			PresentingComplaints: stringPtr("New complaint"),
		}

		result, err := service.UpdateNote(1, input, 3) // Different doctor ID

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "only the doctor who created the clinical note can make changes", err.Error())
		mockNoteRepo.AssertExpectations(t)
		mockNoteRepo.AssertNotCalled(t, "Update")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		existingNote := &models.ClinicalNote{
			Model:    gorm.Model{ID: 1},
			DoctorID: 2,
		}

		newComplaint := "Updated complaint"
		input := models.UpdateNoteInput{
			PresentingComplaints: &newComplaint,
		}

		mockNoteRepo.On("FindByID", uint(1)).Return(existingNote, nil)
		mockNoteRepo.On("Update", mock.AnythingOfType("*models.ClinicalNote")).Return(errors.New("update failed"))

		result, err := service.UpdateNote(1, input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "update failed", err.Error())
		mockNoteRepo.AssertExpectations(t)
	})
}

func TestDeleteNote(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		existingNote := &models.ClinicalNote{
			Model:    gorm.Model{ID: 1},
			DoctorID: 2,
		}

		mockNoteRepo.On("FindByID", uint(1)).Return(existingNote, nil)
		mockNoteRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeleteNote(1, 2)

		assert.NoError(t, err)
		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("NoteNotFound", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		mockNoteRepo.On("FindByID", uint(1)).Return(&models.ClinicalNote{}, errors.New("not found"))

		err := service.DeleteNote(1, 2)

		assert.Error(t, err)
		assert.Equal(t, "clinical note not found", err.Error())
		mockNoteRepo.AssertExpectations(t)
		mockNoteRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("UnauthorizedDelete", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		existingNote := &models.ClinicalNote{
			Model:    gorm.Model{ID: 1},
			DoctorID: 2,
		}

		mockNoteRepo.On("FindByID", uint(1)).Return(existingNote, nil)

		err := service.DeleteNote(1, 3) // Different doctor ID

		assert.Error(t, err)
		assert.Equal(t, "only the doctor who created the clinical note can delete it", err.Error())
		mockNoteRepo.AssertExpectations(t)
		mockNoteRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockNoteRepo := new(mocks.ClinicalNoteRepository)
		mockAppointmentRepo := new(mocks.AppointmentRepository)
		mockPatientRepo := new(mocks.PatientRepository)

		service := NewClinicalNoteService(mockNoteRepo, mockAppointmentRepo, mockPatientRepo)

		existingNote := &models.ClinicalNote{
			Model:    gorm.Model{ID: 1},
			DoctorID: 2,
		}

		mockNoteRepo.On("FindByID", uint(1)).Return(existingNote, nil)
		mockNoteRepo.On("Delete", uint(1)).Return(errors.New("delete failed"))

		err := service.DeleteNote(1, 2)

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
		mockNoteRepo.AssertExpectations(t)
	})
}

// func stringPtr(s string) *string {
// 	return &s
// }
