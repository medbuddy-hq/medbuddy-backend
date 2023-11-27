package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID           primitive.ObjectID `bson:"_id"`
	Time         time.Time          `bson:"time"`
	Status       string             `bson:"status"`
	MedicationID primitive.ObjectID `bson:"medication_id"`
}

type LatestTaskResponse struct {
	ID           primitive.ObjectID  `bson:"_id"`
	Time         time.Time           `bson:"time"`
	Status       string              `bson:"status"`
	MedicationID primitive.ObjectID  `bson:"medication_id"`
	Medication   MedicationForDosage `bson:"medication"`
}
