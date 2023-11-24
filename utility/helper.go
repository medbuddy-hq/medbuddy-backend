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

func RequestsToPractitionerResponse(pract *model.Practitioner, user *model.User) model.PractitionerResponse {
	return model.PractitionerResponse{
		ID:        pract.ID,
		FullName:  pract.FullName,
		UserID:    pract.UserId,
		Email:     pract.Email,
		User:      *user,
		Expertise: pract.Expertise,
		Title:     pract.Title,
	}
}

func MedicineRequestToMedicine(medicine *model.MedicineRequest) model.Medicine {
	return model.Medicine{
		ID:           medicine.ID,
		Name:         medicine.Name,
		Manufacturer: medicine.Manufacturer,
		Category:     medicine.Category,
		Form:         medicine.Form,
		Strength:     medicine.Strength,
		Dosage:       medicine.Dosage,
		CreatedAt:    medicine.CreatedAt,
		UpdatedAt:    medicine.UpdatedAt,
	}
}
