package mocks

import (
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/stretchr/testify/mock"
)

type AppointmentRepository struct {
	mock.Mock
}

func (m *AppointmentRepository) Create(appointment *models.Appointment) error {
	args := m.Called(appointment)
	return args.Error(0)
}

func (m *AppointmentRepository) FindAll(filters map[string]interface{}) ([]models.Appointment, error) {
	args := m.Called(filters)
	return args.Get(0).([]models.Appointment), args.Error(1)
}

func (m *AppointmentRepository) FindByID(id uint) (*models.Appointment, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Appointment), args.Error(1)
}

func (m *AppointmentRepository) Update(appointment *models.Appointment) error {
	args := m.Called(appointment)
	return args.Error(0)
}

func (m *AppointmentRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
