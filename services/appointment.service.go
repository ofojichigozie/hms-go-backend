package services

import (
	"errors"
	"time"

	"github.com/ofojichigozie/hms-go-backend/constants"
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/repositories"
)

type AppointmentService interface {
	CreateAppointment(input models.CreateAppointmentInput, createdBy uint) (*models.Appointment, error)
	GetAllAppointments(filters map[string]interface{}) ([]models.Appointment, error)
	GetAppointmentByID(id uint) (*models.Appointment, error)
	UpdateAppointment(id uint, input models.UpdateAppointmentInput, updatedBy uint) (*models.Appointment, error)
	DeleteAppointment(id uint) error
}

type appointmentService struct {
	appointmentRepository repositories.AppointmentRepository
	patientRepository     repositories.PatientRepository
}

func NewAppointmentService(
	appointmentRepository repositories.AppointmentRepository,
	patientRepository repositories.PatientRepository,
) AppointmentService {
	return &appointmentService{
		appointmentRepository: appointmentRepository,
		patientRepository:     patientRepository,
	}
}

func (as *appointmentService) CreateAppointment(input models.CreateAppointmentInput, createdBy uint) (*models.Appointment, error) {
	patient, err := as.patientRepository.FindByID(input.PatientID)
	if err != nil || patient == nil {
		return nil, errors.New("associated patient record not found")
	}

	appointment := &models.Appointment{
		PatientID:      input.PatientID,
		ReceptionistID: createdBy,
		Department:     input.Department,
		ScheduledAt:    time.Now(),
		Duration:       input.Duration,
		Reason:         input.Reason,
		Status:         constants.AppointmentStatus.SCHEDULED,
	}

	if err := as.appointmentRepository.Create(appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

func (as *appointmentService) GetAllAppointments(filters map[string]interface{}) ([]models.Appointment, error) {
	return as.appointmentRepository.FindAll(filters)
}

func (as *appointmentService) GetAppointmentByID(id uint) (*models.Appointment, error) {
	return as.appointmentRepository.FindByID(id)
}

func (as *appointmentService) UpdateAppointment(id uint, input models.UpdateAppointmentInput, updatedBy uint) (*models.Appointment, error) {
	appointment, err := as.appointmentRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if input.DoctorID != nil {
		appointment.DoctorID = input.DoctorID
	}
	if input.Department != nil {
		appointment.Department = *input.Department
	}
	if input.Status != nil {
		appointment.Status = *input.Status
	}
	if input.Reason != nil {
		appointment.Reason = *input.Reason
	}

	appointment.UpdatedBy = updatedBy

	if err := as.appointmentRepository.Update(appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

func (as *appointmentService) DeleteAppointment(id uint) error {
	appointment, err := as.appointmentRepository.FindByID(id)
	if err != nil {
		return err
	}
	if appointment == nil {
		return errors.New("appointment not found")
	}
	return as.appointmentRepository.Delete(id)
}
