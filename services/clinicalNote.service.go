package services

import (
	"errors"

	"github.com/ofojichigozie/hms-go-backend/constants"
	"github.com/ofojichigozie/hms-go-backend/models"
	"github.com/ofojichigozie/hms-go-backend/repositories"
)

type ClinicalNoteService interface {
	CreateNote(input models.CreateNoteInput, doctorID uint) (*models.ClinicalNote, error)
	GetNoteByID(id uint) (*models.ClinicalNote, error)
	GetNoteByAppointmentID(appointmentID uint) (*models.ClinicalNote, error)
	GetNotesByPatientID(patientID uint) ([]models.ClinicalNote, error)
	UpdateNote(id uint, input models.UpdateNoteInput, staffId uint) (*models.ClinicalNote, error)
	DeleteNote(id uint, staffId uint) error
}

type clinicalNoteService struct {
	clinicalNoteRepository repositories.ClinicalNoteRepository
	appointmentRespository repositories.AppointmentRepository
	patientRespository     repositories.PatientRepository
}

func NewClinicalNoteService(
	clinicalNoteRepository repositories.ClinicalNoteRepository,
	appointmentRespository repositories.AppointmentRepository,
	patientRespository repositories.PatientRepository,
) ClinicalNoteService {
	return &clinicalNoteService{
		clinicalNoteRepository: clinicalNoteRepository,
		appointmentRespository: appointmentRespository,
		patientRespository:     patientRespository,
	}
}

func (cns *clinicalNoteService) CreateNote(input models.CreateNoteInput, doctorID uint) (*models.ClinicalNote, error) {
	appointment, err := cns.appointmentRespository.FindByID(input.AppointmentID)
	if err != nil {
		return nil, errors.New("associated appointment record not found")
	}

	clinicalNote := &models.ClinicalNote{
		AppointmentID:        input.AppointmentID,
		PatientID:            appointment.PatientID,
		DoctorID:             doctorID,
		PresentingComplaints: input.PresentingComplaints,
		PastMedicalHistory:   input.PastMedicalHistory,
		ClinicalDiagnosis:    input.ClinicalDiagnosis,
		TreatmentPlan:        input.TreatmentPlan,
		Recommendation:       input.Recommendation,
	}

	if err := cns.clinicalNoteRepository.Create(clinicalNote); err != nil {
		return nil, err
	}

	appointment.Status = constants.AppointmentStatus.COMPLETED
	cns.appointmentRespository.Update(appointment)

	return clinicalNote, nil
}

func (cns *clinicalNoteService) GetNoteByID(id uint) (*models.ClinicalNote, error) {
	return cns.clinicalNoteRepository.FindByID(id)
}

func (cns *clinicalNoteService) GetNoteByAppointmentID(appointmentID uint) (*models.ClinicalNote, error) {
	return cns.clinicalNoteRepository.FindByAppointmentID(appointmentID)
}

func (cns *clinicalNoteService) GetNotesByPatientID(patientID uint) ([]models.ClinicalNote, error) {
	patient, err := cns.patientRespository.FindByID(patientID)
	if err != nil || patient == nil {
		return nil, errors.New("patient record not found")
	}

	return cns.clinicalNoteRepository.FindByPatientID(patientID)
}

func (cns *clinicalNoteService) UpdateNote(id uint, input models.UpdateNoteInput, staffId uint) (*models.ClinicalNote, error) {
	clinicalNote, err := cns.clinicalNoteRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("clinical note not found")
	}
	if clinicalNote.DoctorID != staffId {
		return nil, errors.New("only the doctor who created the clinical note can make changes")
	}

	if input.PresentingComplaints != nil {
		clinicalNote.PresentingComplaints = *input.PresentingComplaints
	}
	if input.PastMedicalHistory != nil {
		clinicalNote.PastMedicalHistory = *input.PastMedicalHistory
	}
	if input.ClinicalDiagnosis != nil {
		clinicalNote.ClinicalDiagnosis = *input.ClinicalDiagnosis
	}
	if input.TreatmentPlan != nil {
		clinicalNote.TreatmentPlan = *input.TreatmentPlan
	}
	if input.Recommendation != nil {
		clinicalNote.Recommendation = *input.Recommendation
	}

	if err := cns.clinicalNoteRepository.Update(clinicalNote); err != nil {
		return nil, err
	}

	return clinicalNote, nil
}

func (cns *clinicalNoteService) DeleteNote(id uint, staffId uint) error {
	clinicalNote, err := cns.clinicalNoteRepository.FindByID(id)
	if err != nil {
		return errors.New("clinical note not found")
	}
	if clinicalNote.DoctorID != staffId {
		return errors.New("only the doctor who created the clinical note can delete it")
	}

	return cns.clinicalNoteRepository.Delete(id)
}
