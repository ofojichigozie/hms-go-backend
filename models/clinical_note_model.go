package models

import "gorm.io/gorm"

type ClinicalNote struct {
	gorm.Model
	AppointmentID        uint   `json:"appointmentId" gorm:"not null;unique"`
	PatientID            uint   `json:"patientId" gorm:"not null"`
	DoctorID             uint   `json:"doctorId" gorm:"not null"`
	PresentingComplaints string `json:"presentingComplaints" gorm:"type:text"`
	PastMedicalHistory   string `json:"pastMedicalHistory" gorm:"type:text"`
	ClinicalDiagnosis    string `json:"clinicalHistoryDiagnosis" gorm:"type:text"`
	TreatmentPlan        string `json:"treatmentPlan" gorm:"type:text"`
	Recommendation       string `json:"recommendation" gorm:"type:text"`
}

type CreateNoteInput struct {
	AppointmentID        uint   `json:"appointmentId" binding:"required"`
	PresentingComplaints string `json:"presentingComplaints" binding:"required,max=1000"`
	PastMedicalHistory   string `json:"pastMedicalHistory" binding:"max=1000"`
	ClinicalDiagnosis    string `json:"clinicalDiagnosis" binding:"max=1000"`
	TreatmentPlan        string `json:"treatmentPlan" binding:"required,max=1000"`
	Recommendation       string `json:"recommendation" binding:"required,max=1000"`
}

type UpdateNoteInput struct {
	PresentingComplaints *string `json:"presentingComplaints,omitempty" binding:"omitempty,max=1000"`
	PastMedicalHistory   *string `json:"pastMedicalHistory,omitempty" binding:"omitempty,max=1000"`
	ClinicalDiagnosis    *string `json:"clinicalDiagnosis,omitempty" binding:"omitempty,max=1000"`
	TreatmentPlan        *string `json:"treatmentPlan,omitempty" binding:"omitempty,max=1000"`
	Recommendation       *string `json:"recommendation,omitempty" binding:"omitempty,max=1000"`
}
