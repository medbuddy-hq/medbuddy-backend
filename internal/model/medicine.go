package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Medicine struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Name         string             `json:"name,omitempty" bson:"name"`
	Manufacturer string             `json:"manufacturer,omitempty" bson:"manufacturer"`
	Category     string             `json:"category,omitempty" bson:"category"`
	Form         string             `json:"form,omitempty" bson:"form"`
	Strength     string             `json:"strength,omitempty" bson:"strength"`
	Dosage       string             `json:"dosage,omitempty" bson:"dosage"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type MedicineRequest struct {
	ID           primitive.ObjectID `json:"_id"`
	Name         string             `json:"name,omitempty" validate:"required"`
	Manufacturer string             `json:"manufacturer,omitempty" validate:"required"`
	Category     string             `json:"category,omitempty"`
	Form         string             `json:"form,omitempty" validate:"required,oneof='Capsule' 'Solution' 'Tablet' 'Others'"`
	Strength     string             `json:"strength,omitempty" validate:"required"`
	Dosage       string             `json:"dosage,omitempty" validate:"required"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type MedicineFilter struct {
	Name         string `json:"name" validate:"required"`
	Manufacturer string `json:"manufacturer" validate:"required"`
	Strength     string `json:"strength" validate:"required"`
	Form         string `json:"form"`
}
