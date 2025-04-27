package services

import (
	"errors"
	"time"

	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/repositories"
	"github.com/ofojichigozie/hms-go-backend/utils"
)

type PatientService interface {
	CreatePatient(input models.CreatePatientInput, createdBy uint) (*models.Patient, error)
	GetAllPatients(filters map[string]interface{}) ([]models.Patient, error)
	GetPatientByID(id uint) (*models.Patient, error)
	GetPatientByRegistrationNumber(regNumber string) (*models.Patient, error)
	UpdatePatient(id uint, input models.UpdatePatientInput, updatedBy uint) (*models.Patient, error)
	DeletePatient(id uint) error
}

type patientService struct {
	patientRepository repositories.PatientRepository
	staffRepository   repositories.StaffRepository
}

func NewPatientService(patientRepository repositories.PatientRepository,
	staffRepository repositories.StaffRepository) PatientService {
	return &patientService{
		patientRepository: patientRepository,
		staffRepository:   staffRepository,
	}
}

func (ps *patientService) CreatePatient(input models.CreatePatientInput, createdBy uint) (*models.Patient, error) {
	dateOfBirth, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	regNumber := utils.GeneratePatientID()

	patient := &models.Patient{
		RegistrationNumber: regNumber,
		FirstName:          input.FirstName,
		LastName:           input.LastName,
		DateOfBirth:        dateOfBirth,
		Gender:             input.Gender,
		PhoneNumber:        input.PhoneNumber,
		Email:              input.Email,
		Address:            input.Address,
		CreatedBy:          createdBy,
		UpdatedBy:          createdBy,
	}

	if err := ps.patientRepository.Create(patient); err != nil {
		return nil, err
	}

	return patient, nil
}

func (ps *patientService) GetAllPatients(filters map[string]interface{}) ([]models.Patient, error) {
	return ps.patientRepository.FindAll(filters)
}

func (ps *patientService) GetPatientByID(id uint) (*models.Patient, error) {
	return ps.patientRepository.FindByID(id)
}

func (ps *patientService) GetPatientByRegistrationNumber(regNumber string) (*models.Patient, error) {
	return ps.patientRepository.FindByRegistrationNumber(regNumber)
}

func (ps *patientService) UpdatePatient(id uint, input models.UpdatePatientInput, updatedBy uint) (*models.Patient, error) {
	patient, err := ps.patientRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("patient not found")
	}

	if input.FirstName != nil {
		patient.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		patient.LastName = *input.LastName
	}
	if input.DateOfBirth != nil {
		dateOfBirth, err := time.Parse("2006-01-02", *input.DateOfBirth)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		patient.DateOfBirth = dateOfBirth
	}
	if input.Gender != nil {
		patient.Gender = *input.Gender
	}
	if input.PhoneNumber != nil {
		patient.PhoneNumber = *input.PhoneNumber
	}
	if input.Email != nil {
		patient.Email = *input.Email
	}
	if input.Address != nil {
		patient.Address = *input.Address
	}
	if input.BloodGroup != nil {
		patient.BloodGroup = *input.BloodGroup
	}
	if input.Genotype != nil {
		patient.Genotype = *input.Genotype
	}
	patient.UpdatedBy = updatedBy

	if err := ps.patientRepository.Update(patient); err != nil {
		return nil, err
	}

	return patient, nil
}

func (ps *patientService) DeletePatient(id uint) error {
	patient, err := ps.patientRepository.FindByID(id)
	if err != nil {
		return err
	}
	if patient == nil {
		return errors.New("patient not found")
	}
	return ps.patientRepository.Delete(id)
}
