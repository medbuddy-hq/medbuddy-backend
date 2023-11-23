package utility

import "medbuddy-backend/internal/model"

func RequestsToPatientResponse(patient *model.Patient, user *model.User) model.PatientResponse {
	return model.PatientResponse{
		ID:       patient.ID,
		FullName: patient.FullName,
		UserID:   patient.UserID,
		Email:    patient.Email,
		User:     *user,
	}
}
