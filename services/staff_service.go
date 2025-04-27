package services

import (
	"errors"
	"strings"

	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/repositories"
	"github.com/ofojichigozie/hms-go-backend/utils"
)

type StaffService interface {
	CreateStaff(input models.CreateStaffInput) (*models.Staff, error)
	GetAllStaff() ([]models.Staff, error)
	GetStaffByID(id uint) (*models.Staff, error)
	GetStaffByEmail(email string) (*models.Staff, error)
	GetStaffByEmployeeID(employeeID string) (*models.Staff, error)
	UpdateStaff(id uint, input models.UpdateStaffInput) (*models.Staff, error)
	DeleteStaff(id uint) error
}

type staffService struct {
	staffRepository repositories.StaffRepository
}

func NewStaffService(staffRepository repositories.StaffRepository) StaffService {
	return &staffService{staffRepository: staffRepository}
}

func (ss *staffService) CreateStaff(input models.CreateStaffInput) (*models.Staff, error) {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	staff := &models.Staff{
		EmployeeID:     input.EmployeeID,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		PhoneNumber:    input.PhoneNumber,
		Email:          strings.ToLower(input.Email),
		PasswordHash:   hashedPassword,
		Role:           input.Role,
		LicenseNumber:  input.LicenseNumber,
		Specialization: input.Specialization,
		Department:     input.Department,
	}

	if err := ss.staffRepository.Create(staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (ss *staffService) GetAllStaff() ([]models.Staff, error) {
	return ss.staffRepository.FindAll()
}

func (ss *staffService) GetStaffByID(id uint) (*models.Staff, error) {
	staff, err := ss.staffRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errors.New("staff not found")
	}
	return staff, nil
}

func (ss *staffService) GetStaffByEmail(email string) (*models.Staff, error) {
	staff, err := ss.staffRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errors.New("staff not found")
	}
	return staff, nil
}

func (ss *staffService) GetStaffByEmployeeID(employeeID string) (*models.Staff, error) {
	staff, err := ss.staffRepository.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errors.New("staff not found")
	}
	return staff, nil
}

func (ss *staffService) UpdateStaff(id uint, input models.UpdateStaffInput) (*models.Staff, error) {
	staff, err := ss.staffRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, errors.New("staff not found")
	}

	if input.FirstName != nil {
		staff.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		staff.LastName = *input.LastName
	}
	if input.Email != nil {
		staff.Email = strings.ToLower(*input.Email)
	}
	if input.PhoneNumber != nil {
		staff.PhoneNumber = *input.PhoneNumber
	}
	if input.IsActive != nil {
		staff.IsActive = *input.IsActive
	}
	if input.LicenseNumber != nil {
		staff.LicenseNumber = input.LicenseNumber
	}
	if input.Specialization != nil {
		staff.Specialization = input.Specialization
	}
	if input.Department != nil {
		staff.Department = input.Department
	}

	if err := ss.staffRepository.Update(staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (ss *staffService) DeleteStaff(id uint) error {
	staff, err := ss.staffRepository.FindByID(id)
	if err != nil {
		return err
	}
	if staff == nil {
		return errors.New("staff not found")
	}
	return ss.staffRepository.Delete(id)
}
