package mocks

import (
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/stretchr/testify/mock"
)

type PatientRepository struct {
	mock.Mock
}

func (m *PatientRepository) Create(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *PatientRepository) FindAll(filters map[string]interface{}) ([]models.Patient, error) {
	args := m.Called(filters)
	return args.Get(0).([]models.Patient), args.Error(1)
}

func (m *PatientRepository) FindByID(id uint) (*models.Patient, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *PatientRepository) FindByRegistrationNumber(regNumber string) (*models.Patient, error) {
	args := m.Called(regNumber)
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *PatientRepository) Update(patient *models.Patient) error {
	args := m.Called(patient)
	return args.Error(0)
}

func (m *PatientRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
