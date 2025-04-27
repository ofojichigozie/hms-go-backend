package repositories

import (
	"github.com/ofojichigozie/hms-go-backend/models"
	"gorm.io/gorm"
)

type PatientRepository interface {
	Create(patient *models.Patient) error
	FindAll(filters map[string]interface{}) ([]models.Patient, error)
	FindByID(id uint) (*models.Patient, error)
	FindByRegistrationNumber(regNumber string) (*models.Patient, error)
	Update(patient *models.Patient) error
	Delete(id uint) error
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (pr *patientRepository) Create(patient *models.Patient) error {
	return pr.db.Create(patient).Error
}

func (pr *patientRepository) FindAll(filters map[string]interface{}) ([]models.Patient, error) {
	var patients []models.Patient
	query := pr.db.Model(&models.Patient{})

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Find(&patients).Error
	return patients, err
}

func (pr *patientRepository) FindByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	err := pr.db.First(&patient, id).Error
	return &patient, err
}

func (pr *patientRepository) FindByRegistrationNumber(regNumber string) (*models.Patient, error) {
	var patient models.Patient
	err := pr.db.Where("registration_number = ?", regNumber).First(&patient).Error
	return &patient, err
}

func (pr *patientRepository) Update(patient *models.Patient) error {
	return pr.db.Save(patient).Error
}

func (pr *patientRepository) Delete(id uint) error {
	return pr.db.Delete(&models.Patient{}, id).Error
}
