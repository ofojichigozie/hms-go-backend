package repositories

import (
	"github.com/ofojichigozie/hms-go-backend/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	Create(appointment *models.Appointment) error
	FindAll(filters map[string]interface{}) ([]models.Appointment, error)
	FindByID(id uint) (*models.Appointment, error)
	Update(appointment *models.Appointment) error
	Delete(id uint) error
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (ar *appointmentRepository) Create(appointment *models.Appointment) error {
	return ar.db.Create(appointment).Error
}

func (ar *appointmentRepository) FindAll(filters map[string]interface{}) ([]models.Appointment, error) {
	var appointments []models.Appointment
	query := ar.db.Model(&models.Appointment{})

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Find(&appointments).Error
	return appointments, err
}

func (ar *appointmentRepository) FindByID(id uint) (*models.Appointment, error) {
	var appointment models.Appointment
	err := ar.db.First(&appointment, id).Error
	return &appointment, err
}

func (ar *appointmentRepository) Update(appointment *models.Appointment) error {
	return ar.db.Save(appointment).Error
}

func (ar *appointmentRepository) Delete(id uint) error {
	return ar.db.Delete(&models.Appointment{}, id).Error
}
