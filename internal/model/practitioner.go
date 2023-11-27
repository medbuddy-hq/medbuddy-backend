package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Practitioner struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FullName  string             `bson:"full_name,omitempty"`
	Title     string             `bson:"title,omitempty"`
	Email     string             `bson:"email,omitempty"`
	UserId    primitive.ObjectID `bson:"user_id,omitempty"`
	Expertise string             `bson:"expertise,omitempty"`
}

type PractitionerRequest struct {
	Firstname string `json:"firstname,omitempty" validate:"required"`
	Lastname  string `json:"lastname,omitempty" validate:"required"`
	DOB       string `json:"dob,omitempty"`
	Gender    string `json:"gender" validate:"required,oneof='male' 'female''"`
	Email     string `json:"email,omitempty" validate:"email,required"`
	Title     string `json:"title,omitempty"`
	Expertise string `json:"expertise,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required,min=8"`
}

type PractitionerResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	FullName  string             `json:"fullname,omitempty" bson:"full_name"`
	Email     string             `json:"email,omitempty" bson:"email"`
	Expertise string             `json:"expertise,omitempty"`
	Title     string             `json:"title,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id"`
	User      User               `json:"user" bson:"user"`
	Token     string             `json:"token,omitempty"`
}
