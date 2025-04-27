package mocks

import (
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/stretchr/testify/mock"
)

type StaffRepository struct {
	mock.Mock
}

func (m *StaffRepository) Create(staff *models.Staff) error {
	args := m.Called(staff)
	return args.Error(0)
}

func (m *StaffRepository) FindAll() ([]models.Staff, error) {
	args := m.Called()
	return args.Get(0).([]models.Staff), args.Error(1)
}

func (m *StaffRepository) FindByID(id uint) (*models.Staff, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *StaffRepository) FindByEmail(email string) (*models.Staff, error) {
	args := m.Called(email)
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *StaffRepository) FindByEmployeeID(employeeID string) (*models.Staff, error) {
	args := m.Called(employeeID)
	return args.Get(0).(*models.Staff), args.Error(1)
}

func (m *StaffRepository) Update(staff *models.Staff) error {
	args := m.Called(staff)
	return args.Error(0)
}

func (m *StaffRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
