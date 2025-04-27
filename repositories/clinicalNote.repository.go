package repositories

import (
	"github.com/ofojichigozie/hms-go-backend/models"
	"gorm.io/gorm"
)

type ClinicalNoteRepository interface {
	Create(note *models.ClinicalNote) error
	FindByID(id uint) (*models.ClinicalNote, error)
	FindByAppointmentID(appointmentID uint) (*models.ClinicalNote, error)
	FindByPatientID(patientID uint) ([]models.ClinicalNote, error)
	Update(note *models.ClinicalNote) error
	Delete(id uint) error
}

type clinicalNoteRepository struct {
	db *gorm.DB
}

func NewClinicalNoteRepository(db *gorm.DB) ClinicalNoteRepository {
	return &clinicalNoteRepository{db}
}

func (r *clinicalNoteRepository) Create(note *models.ClinicalNote) error {
	return r.db.Create(note).Error
}

func (r *clinicalNoteRepository) FindByID(id uint) (*models.ClinicalNote, error) {
	var note models.ClinicalNote
	err := r.db.First(&note, id).Error
	return &note, err
}

func (r *clinicalNoteRepository) FindByAppointmentID(appointmentID uint) (*models.ClinicalNote, error) {
	var note models.ClinicalNote
	err := r.db.Where("appointment_id = ?", appointmentID).First(&note).Error
	return &note, err
}

func (r *clinicalNoteRepository) FindByPatientID(patientID uint) ([]models.ClinicalNote, error) {
	var notes []models.ClinicalNote
	err := r.db.Where("patient_id = ?", patientID).Find(&notes).Error
	return notes, err
}

func (r *clinicalNoteRepository) Update(note *models.ClinicalNote) error {
	return r.db.Save(note).Error
}

func (r *clinicalNoteRepository) Delete(id uint) error {
	return r.db.Delete(&models.ClinicalNote{}, id).Error
}
