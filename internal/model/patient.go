package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Patient struct {
	ID       primitive.ObjectID `bson:"_id"`
	FullName string             `bson:"full_name"`
	Email    string             `bson:"email"`
	UserID   primitive.ObjectID `bson:"user_id"`
}

type CreatePatientReq struct {
	Firstname string `json:"firstname,omitempty" validate:"required"`
	Lastname  string `json:"lastname,omitempty" validate:"required"`
	DOB       string `json:"dob,omitempty" validate:"required"`
	Gender    string `json:"gender" validate:"required,oneof='male' 'female''"`
	Email     string `json:"email,omitempty" validate:"email,required"`
	Password  string `json:"password,omitempty" validate:"required,min=8"`
}

type PatientResponse struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	FullName string             `json:"fullname,omitempty" bson:"full_name"`
	Email    string             `json:"email,omitempty" bson:"email"`
	UserID   primitive.ObjectID `json:"user_id,omitempty" bson:"user_id"`
	User     User               `json:"user" bson:"user"`
	Token    string             `json:"token,omitempty"`
}
