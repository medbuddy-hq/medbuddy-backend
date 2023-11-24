package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Dosage struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	ReminderTime time.Time          `json:"reminder_time" bson:"reminder_time"`
	Status       string             `json:"status" bson:"status"`
	TimeTaken    time.Time          `json:"time_taken,omitempty" bson:"time_taken"`
	TimeSkipped  time.Time          `json:"time_skipped,omitempty" bson:"time_skipped"`
	IsActive     bool               `json:"is_active" bson:"is_active"`
	MedicationID primitive.ObjectID `json:"medication_id" bson:"medication_id"`
	PatientID    primitive.ObjectID `json:"patient_id" bson:"patient_id"`
}

type DosageResponse struct {
	ID           primitive.ObjectID  `json:"_id" bson:"_id"`
	ReminderTime time.Time           `json:"reminder_time" bson:"reminder_time"`
	Status       string              `json:"status" bson:"status"`
	TimeTaken    time.Time           `json:"time_taken,omitempty" bson:"time_taken"`
	TimeSkipped  time.Time           `json:"time_skipped,omitempty" bson:"time_skipped"`
	IsActive     bool                `json:"is_active" bson:"is_active"`
	MedicationID primitive.ObjectID  `json:"medication_id" bson:"medication_id"`
	Medication   MedicationForDosage `json:"medication" bson:"medication"`
	PatientID    primitive.ObjectID  `json:"patient_id" bson:"patient_id"`
}

type MedicationForDosage struct {
	Name                string    `bson:"name" json:"name,omitempty"`
	StartDate           time.Time `bson:"start_date" json:"start_date"`                                   // date to start taking the medicine
	DosageQuantity      string    `bson:"dosage_quantity" json:"dosage_quantity,omitempty"`               // measure (quantity) of medicine taken per dosage
	DailyDosage         int       `bson:"daily_dosage" json:"daily_dosage,omitempty"`                     // measure (quantity) of dosage per day
	TotalNumberOfDosage int       `bson:"total_number_of_dosage" json:"total_number_of_dosage,omitempty"` // total number of dosages
	DosagesTaken        int       `bson:"dosages_taken" json:"dosages_taken,omitempty"`
	Treatment           string    `bson:"treatment" json:"treatment,omitempty"` // sickness/disease
	Comment             string    `bson:"comment" json:"comment,omitempty"`
	Medicine            Medicine  `bson:"medicine" json:"medicine"`
}

type DosageFilter struct {
	PatiendID    primitive.ObjectID
	MedicationID primitive.ObjectID
	IsActive     *bool
}

type SetStatusRequest struct {
	Status string `json:"status" validate:"required"`
}
