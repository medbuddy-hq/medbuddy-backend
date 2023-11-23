package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Medication struct {
	ID                  primitive.ObjectID `bson:"_id"`
	Name                string             `bson:"name"`
	StartDate           time.Time          `bson:"start_date"`
	EndDate             time.Time          `bson:"end_date"`
	DosageQuantity      string             `bson:"dosage_quantity"`
	DailyDosage         int                `bson:"daily_dosage"`
	TotalNumberOfDosage int                `bson:"total_number_of_dosage"`
	DosageTimes         []time.Time        `bson:"dosage_times"`
	Treatment           string             `bson:"treatment"`
	CreatedAt           time.Time          `bson:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at"`
	PatientID           primitive.ObjectID `bson:"patient_id"`
	MedicineID          primitive.ObjectID `bson:"medicine_id"`
}

type MedicationRequest struct {
	Name           string             `json:"name,omitempty"`
	StartDate      time.Time          `json:"start_date"`
	EndDate        time.Time          `json:"end_date"`
	DosageQuantity string             `json:"dosage_quantity,omitempty"`
	DailyDosage    int                `json:"daily_dosage,omitempty"`
	DosageTimes    []time.Time        `json:"dosage_times,omitempty"`
	Treatment      string             `json:"treatment,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	MedicineID     primitive.ObjectID `json:"medicine_id,omitempty"`
	Medicine       Medicine           `json:"medicine"`
}

type MedicationResponse struct {
	ID             primitive.ObjectID `json:"id,omitempty"`
	Name           string             `json:"name,omitempty"`
	StartDate      time.Time          `json:"start_date"`
	EndDate        time.Time          `json:"end_date"`
	DosageQuantity string             `json:"dosage_quantity,omitempty"`
	DailyDosage    int                `json:"daily_dosage,omitempty"`
	DosageTimes    []time.Time        `json:"dosage_times,omitempty"`
	Treatment      string             `json:"treatment,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	MedicineID     primitive.ObjectID `json:"medicine_id,omitempty"`
	Medicine       Medicine           `json:"medicine"`
	PatientID      primitive.ObjectID `json:"patient_id"`
	Patient        Patient            `json:"patient"`
}
