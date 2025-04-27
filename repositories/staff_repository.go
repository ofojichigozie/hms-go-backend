package repositories

import (
	"errors"
	"strings"

	"github.com/ofojichigozie/hms-go-backend/models"
	"gorm.io/gorm"
)

type StaffRepository interface {
	Create(staff *models.Staff) error
	FindAll() ([]models.Staff, error)
	FindByID(id uint) (*models.Staff, error)
	FindByEmail(email string) (*models.Staff, error)
	FindByEmployeeID(employeeID string) (*models.Staff, error)
	Update(staff *models.Staff) error
	Delete(id uint) error
}

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{db: db}
}

func (sr *staffRepository) Create(staff *models.Staff) error {
	return sr.db.Create(staff).Error
}

func (sr *staffRepository) FindAll() ([]models.Staff, error) {
	var staff []models.Staff
	err := sr.db.Find(&staff).Error
	return staff, err
}

func (sr *staffRepository) FindByID(id uint) (*models.Staff, error) {
	var staff models.Staff
	err := sr.db.First(&staff, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &staff, err
}

func (sr *staffRepository) FindByEmail(email string) (*models.Staff, error) {
	var staff models.Staff
	err := sr.db.Where("email = ?", strings.ToLower(email)).First(&staff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &staff, err
}

func (sr *staffRepository) FindByEmployeeID(employeeID string) (*models.Staff, error) {
	var staff models.Staff
	err := sr.db.Where("employee_id = ?", employeeID).First(&staff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &staff, err
}

func (sr *staffRepository) Update(staff *models.Staff) error {
	return sr.db.Save(staff).Error
}

func (sr *staffRepository) Delete(id uint) error {
	return sr.db.Delete(&models.Staff{}, id).Error
}
