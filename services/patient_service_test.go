package services

import (
	"errors"
	"testing"
	"time"

	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreatePatient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		input := models.CreatePatientInput{
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: "1990-01-01",
			Gender:      "male",
			PhoneNumber: "+1234567890",
			Email:       "john@example.com",
			Address:     "123 Main St",
		}

		mockPatientRepo.On("Create", mock.AnythingOfType("*models.Patient")).Return(nil).Run(func(args mock.Arguments) {
			patient := args.Get(0).(*models.Patient)
			assert.NotEmpty(t, patient.RegistrationNumber)
			assert.Equal(t, input.FirstName, patient.FirstName)
			assert.Equal(t, input.LastName, patient.LastName)
			assert.Equal(t, "male", patient.Gender)
			assert.Equal(t, input.PhoneNumber, patient.PhoneNumber)
			assert.Equal(t, input.Email, patient.Email)
			assert.Equal(t, input.Address, patient.Address)
			assert.Equal(t, uint(1), patient.CreatedBy)
		})

		result, err := service.CreatePatient(input, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockPatientRepo.AssertExpectations(t)
	})

	t.Run("InvalidDateOfBirth", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		input := models.CreatePatientInput{
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: "invalid-date",
			Gender:      "male",
			PhoneNumber: "+1234567890",
		}

		result, err := service.CreatePatient(input, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid date format, use YYYY-MM-DD", err.Error())
		mockPatientRepo.AssertNotCalled(t, "Create")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		input := models.CreatePatientInput{
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: "1990-01-01",
			Gender:      "male",
			PhoneNumber: "+1234567890",
		}

		mockPatientRepo.On("Create", mock.AnythingOfType("*models.Patient")).Return(errors.New("database error"))

		result, err := service.CreatePatient(input, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockPatientRepo.AssertExpectations(t)
	})
}

func TestGetAllPatients(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		filters := map[string]interface{}{"gender": "male"}
		expectedPatients := []models.Patient{
			{
				Model:              gorm.Model{ID: 1},
				RegistrationNumber: "PAT001",
				FirstName:          "John",
				LastName:           "Doe",
				Gender:             "male",
			},
		}

		mockPatientRepo.On("FindAll", filters).Return(expectedPatients, nil)

		result, err := service.GetAllPatients(filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedPatients, result)
		mockPatientRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		filters := map[string]interface{}{}
		mockPatientRepo.On("FindAll", filters).Return([]models.Patient{}, errors.New("database error"))

		result, err := service.GetAllPatients(filters)

		assert.Error(t, err)
		assert.Equal(t, []models.Patient{}, result)
		assert.Equal(t, "database error", err.Error())
		mockPatientRepo.AssertExpectations(t)
	})
}

func TestGetPatientByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		expectedPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
			FirstName:          "John",
			LastName:           "Doe",
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(expectedPatient, nil)

		result, err := service.GetPatientByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedPatient, result)
		mockPatientRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		mockPatientRepo.On("FindByID", uint(1)).Return(&models.Patient{}, errors.New("not found"))

		result, err := service.GetPatientByID(1)

		assert.Error(t, err)
		assert.Equal(t, &models.Patient{}, result)
		mockPatientRepo.AssertExpectations(t)
	})
}

func TestGetPatientByRegistrationNumber(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		expectedPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
			FirstName:          "John",
			LastName:           "Doe",
		}

		mockPatientRepo.On("FindByRegistrationNumber", "PAT001").Return(expectedPatient, nil)

		result, err := service.GetPatientByRegistrationNumber("PAT001")

		assert.NoError(t, err)
		assert.Equal(t, expectedPatient, result)
		mockPatientRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		mockPatientRepo.On("FindByRegistrationNumber", "PAT001").Return(&models.Patient{}, errors.New("not found"))

		result, err := service.GetPatientByRegistrationNumber("PAT001")

		assert.Error(t, err)
		assert.Equal(t, &models.Patient{}, result)
		mockPatientRepo.AssertExpectations(t)
	})
}

func TestUpdatePatient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		existingPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
			FirstName:          "John",
			LastName:           "Doe",
			DateOfBirth:        time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:             "male",
			PhoneNumber:        "+1234567890",
		}

		newFirstName := "Jonathan"
		newEmail := "jonathan@example.com"
		newDateOfBirth := "1991-02-02"

		input := models.UpdatePatientInput{
			FirstName:   &newFirstName,
			DateOfBirth: &newDateOfBirth,
			Email:       &newEmail,
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(existingPatient, nil)
		mockPatientRepo.On("Update", mock.AnythingOfType("*models.Patient")).Return(nil)

		result, err := service.UpdatePatient(1, input, 2)

		assert.NoError(t, err)
		assert.Equal(t, newFirstName, result.FirstName)
		assert.Equal(t, newEmail, result.Email)
		assert.Equal(t, time.Date(1991, 2, 2, 0, 0, 0, 0, time.UTC), result.DateOfBirth)
		assert.Equal(t, uint(2), result.UpdatedBy)
		mockPatientRepo.AssertExpectations(t)
	})

	t.Run("PatientNotFound", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		mockPatientRepo.On("FindByID", uint(1)).Return(&models.Patient{}, errors.New("not found"))

		input := models.UpdatePatientInput{
			FirstName: stringPtr("Jonathan"),
		}

		result, err := service.UpdatePatient(1, input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "patient not found", err.Error())
		mockPatientRepo.AssertExpectations(t)
		mockPatientRepo.AssertNotCalled(t, "Update")
	})

	t.Run("InvalidDateOfBirth", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		existingPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
			FirstName:          "John",
		}

		invalidDate := "invalid-date"
		input := models.UpdatePatientInput{
			DateOfBirth: &invalidDate,
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(existingPatient, nil)

		result, err := service.UpdatePatient(1, input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid date format, use YYYY-MM-DD", err.Error())
		mockPatientRepo.AssertExpectations(t)
		mockPatientRepo.AssertNotCalled(t, "Update")
	})

	t.Run("RepositoryUpdateError", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		existingPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
			FirstName:          "John",
		}

		newFirstName := "Jonathan"
		input := models.UpdatePatientInput{
			FirstName: &newFirstName,
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(existingPatient, nil)
		mockPatientRepo.On("Update", mock.AnythingOfType("*models.Patient")).Return(errors.New("update failed"))

		result, err := service.UpdatePatient(1, input, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "update failed", err.Error())
		mockPatientRepo.AssertExpectations(t)
	})
}

func TestDeletePatient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		existingPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(existingPatient, nil)
		mockPatientRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeletePatient(1)

		assert.NoError(t, err)
		mockPatientRepo.AssertExpectations(t)
	})

	t.Run("PatientNotFound", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		mockPatientRepo.On("FindByID", uint(1)).Return(&models.Patient{}, errors.New("patient not found"))

		err := service.DeletePatient(1)

		assert.Error(t, err)
		assert.Equal(t, "patient not found", err.Error())
		mockPatientRepo.AssertExpectations(t)
		mockPatientRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("RepositoryDeleteError", func(t *testing.T) {
		mockPatientRepo := new(mocks.PatientRepository)
		mockStaffRepo := new(mocks.StaffRepository)

		service := NewPatientService(mockPatientRepo, mockStaffRepo)

		existingPatient := &models.Patient{
			Model:              gorm.Model{ID: 1},
			RegistrationNumber: "PAT001",
		}

		mockPatientRepo.On("FindByID", uint(1)).Return(existingPatient, nil)
		mockPatientRepo.On("Delete", uint(1)).Return(errors.New("delete failed"))

		err := service.DeletePatient(1)

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
		mockPatientRepo.AssertExpectations(t)
	})
}

// func stringPtr(s string) *string {
// 	return &s
// }
