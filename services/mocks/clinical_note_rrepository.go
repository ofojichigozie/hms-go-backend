package mocks

import (
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/stretchr/testify/mock"
)

type ClinicalNoteRepository struct {
	mock.Mock
}

func (m *ClinicalNoteRepository) Create(note *models.ClinicalNote) error {
	args := m.Called(note)
	return args.Error(0)
}

func (m *ClinicalNoteRepository) FindByID(id uint) (*models.ClinicalNote, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ClinicalNote), args.Error(1)
}

func (m *ClinicalNoteRepository) FindByAppointmentID(appointmentID uint) (*models.ClinicalNote, error) {
	args := m.Called(appointmentID)
	return args.Get(0).(*models.ClinicalNote), args.Error(1)
}

func (m *ClinicalNoteRepository) FindByPatientID(patientID uint) ([]models.ClinicalNote, error) {
	args := m.Called(patientID)
	return args.Get(0).([]models.ClinicalNote), args.Error(1)
}

func (m *ClinicalNoteRepository) Update(note *models.ClinicalNote) error {
	args := m.Called(note)
	return args.Error(0)
}

func (m *ClinicalNoteRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
