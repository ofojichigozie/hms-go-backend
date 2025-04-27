package services

import (
	"errors"
	"strings"
	"testing"

	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateStaff(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		input := models.CreateStaffInput{
			EmployeeID:    "EMP001",
			FirstName:     "John",
			LastName:      "Doe",
			PhoneNumber:   "+1234567890",
			Email:         "john@example.com",
			Password:      "password123",
			Role:          "doctor",
			LicenseNumber: stringPtr("LIC123"),
		}

		mockStaffRepo.On("Create", mock.AnythingOfType("*models.Staff")).Return(nil).Run(func(args mock.Arguments) {
			staff := args.Get(0).(*models.Staff)
			assert.Equal(t, input.EmployeeID, staff.EmployeeID)
			assert.Equal(t, input.FirstName, staff.FirstName)
			assert.Equal(t, input.LastName, staff.LastName)
			assert.Equal(t, input.PhoneNumber, staff.PhoneNumber)
			assert.Equal(t, strings.ToLower(input.Email), staff.Email)
			assert.NotEmpty(t, staff.PasswordHash)
			assert.Equal(t, input.Role, staff.Role)
			assert.Equal(t, *input.LicenseNumber, *staff.LicenseNumber)
		})

		result, err := service.CreateStaff(input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		input := models.CreateStaffInput{
			EmployeeID:    "EMP001",
			FirstName:     "John",
			LastName:      "Doe",
			PhoneNumber:   "+1234567890",
			Email:         "john@example.com",
			Password:      "password123",
			Role:          "doctor",
			LicenseNumber: stringPtr("LIC123"),
		}

		mockStaffRepo.On("Create", mock.AnythingOfType("*models.Staff")).Return(errors.New("database error"))

		result, err := service.CreateStaff(input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("MissingLicenseForDoctor", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		input := models.CreateStaffInput{
			EmployeeID:    "EMP001",
			FirstName:     "John",
			LastName:      "Doe",
			PhoneNumber:   "+1234567890",
			Email:         "john@example.com",
			Password:      "password123",
			Role:          "doctor",
			LicenseNumber: nil,
		}

		mockStaffRepo.On("Create", mock.AnythingOfType("*models.Staff")).Return(errors.New("license number is required for doctors"))

		result, err := service.CreateStaff(input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "license number is required for doctors")
		mockStaffRepo.AssertNotCalled(t, "Create")
	})
}

func TestGetAllStaff(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		expectedStaff := []models.Staff{
			{
				EmployeeID: "EMP001",
				FirstName:  "John",
				LastName:   "Doe",
			},
		}

		mockStaffRepo.On("FindAll").Return(expectedStaff, nil)

		result, err := service.GetAllStaff()

		assert.NoError(t, err)
		assert.Equal(t, expectedStaff, result)
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		mockStaffRepo.On("FindAll").Return([]models.Staff{}, errors.New("database error"))

		result, err := service.GetAllStaff()

		assert.Error(t, err)
		assert.Equal(t, []models.Staff{}, result)
		assert.Equal(t, "database error", err.Error())
		mockStaffRepo.AssertExpectations(t)
	})
}

func TestGetStaffByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		expectedStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
			LastName:   "Doe",
		}

		mockStaffRepo.On("FindByID", uint(1)).Return(expectedStaff, nil)

		result, err := service.GetStaffByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedStaff, result)
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		mockStaffRepo.On("FindByID", uint(1)).Return(&models.Staff{}, errors.New("staff not found"))

		result, err := service.GetStaffByID(1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "staff not found", err.Error())
		mockStaffRepo.AssertExpectations(t)
	})
}

func TestGetStaffByEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		expectedStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john@example.com",
		}

		mockStaffRepo.On("FindByEmail", "john@example.com").Return(expectedStaff, nil)

		result, err := service.GetStaffByEmail("john@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedStaff, result)
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		mockStaffRepo.On("FindByEmail", "john@example.com").Return(&models.Staff{}, errors.New("staff not found"))

		result, err := service.GetStaffByEmail("john@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "staff not found", err.Error())
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("EmailCaseInsensitive", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		expectedStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john@example.com",
		}

		mockStaffRepo.On("FindByEmail", "john@example.com").Return(expectedStaff, nil)

		result, err := service.GetStaffByEmail("JOHN@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedStaff, result)
		mockStaffRepo.AssertExpectations(t)
	})
}

func TestGetStaffByEmployeeID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		expectedStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
			LastName:   "Doe",
		}

		mockStaffRepo.On("FindByEmployeeID", "EMP001").Return(expectedStaff, nil)

		result, err := service.GetStaffByEmployeeID("EMP001")

		assert.NoError(t, err)
		assert.Equal(t, expectedStaff, result)
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		mockStaffRepo.On("FindByEmployeeID", "EMP001").Return(&models.Staff{}, errors.New("staff not found"))

		result, err := service.GetStaffByEmployeeID("EMP001")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "staff not found", err.Error())
		mockStaffRepo.AssertExpectations(t)
	})
}

func TestUpdateStaff(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		existingStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john@example.com",
		}

		newFirstName := "Jonathan"
		newEmail := "jonathan@example.com"
		newIsActive := false

		input := models.UpdateStaffInput{
			FirstName: &newFirstName,
			Email:     &newEmail,
			IsActive:  &newIsActive,
		}

		mockStaffRepo.On("FindByID", uint(1)).Return(existingStaff, nil)
		mockStaffRepo.On("Update", mock.AnythingOfType("*models.Staff")).Return(nil)

		result, err := service.UpdateStaff(1, input)

		assert.NoError(t, err)
		assert.Equal(t, newFirstName, result.FirstName)
		assert.Equal(t, newEmail, result.Email)
		assert.Equal(t, newIsActive, result.IsActive)
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("StaffNotFound", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		mockStaffRepo.On("FindByID", uint(1)).Return(&models.Staff{}, errors.New("staff not found"))

		input := models.UpdateStaffInput{
			FirstName: stringPtr("Jonathan"),
		}

		result, err := service.UpdateStaff(1, input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "staff not found", err.Error())
		mockStaffRepo.AssertExpectations(t)
		mockStaffRepo.AssertNotCalled(t, "Update")
	})

	t.Run("RepositoryUpdateError", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		existingStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
		}

		newFirstName := "Jonathan"
		input := models.UpdateStaffInput{
			FirstName: &newFirstName,
		}

		mockStaffRepo.On("FindByID", uint(1)).Return(existingStaff, nil)
		mockStaffRepo.On("Update", mock.AnythingOfType("*models.Staff")).Return(errors.New("update failed"))

		result, err := service.UpdateStaff(1, input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "update failed", err.Error())
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("EmailCaseInsensitive", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		existingStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
			Email:      "john@example.com",
		}

		newEmail := "JOHN@NEWEXAMPLE.COM"
		input := models.UpdateStaffInput{
			Email: &newEmail,
		}

		mockStaffRepo.On("FindByID", uint(1)).Return(existingStaff, nil)
		mockStaffRepo.On("Update", mock.AnythingOfType("*models.Staff")).Return(nil).Run(func(args mock.Arguments) {
			staff := args.Get(0).(*models.Staff)
			assert.Equal(t, strings.ToLower(newEmail), staff.Email)
		})

		result, err := service.UpdateStaff(1, input)

		assert.NoError(t, err)
		assert.Equal(t, strings.ToLower(newEmail), result.Email)
		mockStaffRepo.AssertExpectations(t)
	})
}

func TestDeleteStaff(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		existingStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
		}

		mockStaffRepo.On("FindByID", uint(1)).Return(existingStaff, nil)
		mockStaffRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeleteStaff(1)

		assert.NoError(t, err)
		mockStaffRepo.AssertExpectations(t)
	})

	t.Run("StaffNotFound", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		mockStaffRepo.On("FindByID", uint(1)).Return(&models.Staff{}, errors.New("staff not found"))

		err := service.DeleteStaff(1)

		assert.Error(t, err)
		assert.Equal(t, "staff not found", err.Error())
		mockStaffRepo.AssertExpectations(t)
		mockStaffRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("RepositoryDeleteError", func(t *testing.T) {
		mockStaffRepo := new(mocks.StaffRepository)
		service := NewStaffService(mockStaffRepo)

		existingStaff := &models.Staff{
			EmployeeID: "EMP001",
			FirstName:  "John",
		}

		mockStaffRepo.On("FindByID", uint(1)).Return(existingStaff, nil)
		mockStaffRepo.On("Delete", uint(1)).Return(errors.New("delete failed"))

		err := service.DeleteStaff(1)

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
		mockStaffRepo.AssertExpectations(t)
	})
}

// func stringPtr(s string) *string {
// 	return &s
// }
