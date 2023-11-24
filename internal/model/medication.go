package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Medication struct {
	ID                  primitive.ObjectID `bson:"_id"`
	Name                string             `bson:"name"`
	StartDate           time.Time          `bson:"start_date"` // date to start taking the medicine
	EndDate             time.Time          `bson:"end_date"`
	DosageQuantity      string             `bson:"dosage_quantity"`        // measure (quantity) of medicine taken per dosage
	DailyDosage         int                `bson:"daily_dosage"`           // measure (quantity) of dosage per day
	TotalNumberOfDosage int                `bson:"total_number_of_dosage"` // total number of dosages
	DosagesTaken        int                `bson:"dosages_taken"`
	Treatment           string             `bson:"treatment"` // sickness/disease
	Comment             string             `bson:"comment"`
	IsActive            bool               `bson:"is_active"`
	CreatedAt           time.Time          `bson:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at"`
	PatientID           primitive.ObjectID `bson:"patient_id"`
	MedicineID          primitive.ObjectID `bson:"medicine_id"`
}

type MedicationRequest struct {
	Name                string    `json:"name,omitempty" validate:"required"`
	StartDate           string    `json:"start_date" validate:"required"`
	EndDate             string    `json:"end_date"`
	DosageQuantity      string    `json:"dosage_quantity,omitempty" validate:"required"`
	DailyDosage         int       `json:"daily_dosage,omitempty" validate:"required"`
	TotalNumberOfDosage int       `json:"total_number_of_dosage"` // total number of dosages
	DosageTimes         []string  `json:"dosage_times,omitempty" validate:"required"`
	Treatment           string    `json:"treatment,omitempty" validate:"required"`
	Comment             string    `json:"comment"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Medicine            Medicine  `json:"medicine" validate:"required"`
}

type MedicationResponse struct {
	ID                  primitive.ObjectID `json:"id,omitempty"`
	Name                string             `json:"name,omitempty"`
	StartDate           time.Time          `json:"start_date"`
	EndDate             time.Time          `json:"end_date"`
	DosageQuantity      string             `json:"dosage_quantity,omitempty"`
	DailyDosage         int                `json:"daily_dosage,omitempty"`
	Dosages             []Dosage           `json:"dosages,omitempty"`
	DosagesTaken        int                `json:"dosages_taken"`
	TotalNumberOfDosage int                `json:"total_number_of_dosage"` // total number of dosages
	Treatment           string             `json:"treatment,omitempty"`
	Comment             string             `json:"comment"`
	IsActive            bool               `json:"is_active"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	MedicineID          primitive.ObjectID `json:"medicine_id,omitempty"`
	Medicine            Medicine           `json:"medicine"`
	PatientID           primitive.ObjectID `json:"patient_id"`
	Patient             Patient            `json:"patient"`
}
