package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medbuddy-backend/internal/model"
)

// repositories

type RedisRepository interface {
	RedisSet(key string, value interface{}) error
	RedisGet(key string) ([]byte, error)
	RedisDelete(key string) (int64, error)
}

// repositories
type StorageRepository interface {
	// Patient
	CreatePatient(ctx context.Context, user *model.Patient) error
	GetPatientByEmail(ctx context.Context, email string) (patient model.PatientResponse, found bool, err error)
	GetPatientByID(ctx context.Context, id primitive.ObjectID) (patient model.PatientResponse, found bool, err error)

	// User
	CreateUser(ctx context.Context, data *model.User) error

	// Medicine
	AddMedicine(ctx context.Context, data *model.Medicine) error
	GetMedicineByID(ctx context.Context, id primitive.ObjectID) (medicine model.Medicine, found bool, err error)
	UpdateMedicine(ctx context.Context, id primitive.ObjectID, data *model.Medicine) (found bool, err error)
	DeleteMedicine(ctx context.Context, id primitive.ObjectID) (found bool, err error)
	GetMedicineFilter(ctx context.Context, req *model.MedicineFilter) (medicine model.Medicine, found bool, err error)

	// Medication
	AddMedication(ctx context.Context, data *model.Medication) error
	UpdateMedication(ctx context.Context, id primitive.ObjectID, data *model.MedicationRequest) (found bool, err error)
	DeleteMedication(ctx context.Context, id primitive.ObjectID) (found bool, err error)
	GetMedication(ctx context.Context, id primitive.ObjectID) (medic model.MedicationResponse, found bool, err error)
	GetPatientsMedications(ctx context.Context, patientId primitive.ObjectID) (medics []model.MedicationResponse, err error)

	// Practitioner

	// Dosage
	SaveDosages(ctx context.Context, data []model.Dosage) error
	GetPatientDosages(ctx context.Context, request *model.DosageFilter) (dosages []model.DosageResponse, err error)
	SetStatus(ctx context.Context, dosageId, patientId primitive.ObjectID, status string) (found bool, err error)
	GetDosage(ctx context.Context, id primitive.ObjectID) (dosage model.DosageResponse, found bool, err error)
	DeleteDosages(ctx context.Context, medicationId primitive.ObjectID) (int64, error)
}
