package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Firstname string             `json:"firstname,omitempty" validate:"required" bson:"firstname"`
	Lastname  string             `json:"lastname,omitempty" validate:"required" bson:"lastname"`
	DOB       time.Time          `json:"dob,omitempty" validate:"required" bson:"dob"`
	Gender    string             `json:"gender" validate:"required,oneof='Male' 'Female''" bson:"gender"`
	Email     string             `json:"email,omitempty" validate:"email,required" bson:"email"`
	Password  string             `json:"password,omitempty" validate:"required,min=8" bson:"password"`
	Role      int                `json:"role,omitempty" bson:"role"`
	IsLocked  bool               `json:"is_locked,omitempty" bson:"is_locked"`
	Salt      string             `json:"salt,omitempty" bson:"salt"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type ResetPassword struct {
	Password string `json:"password,omitempty"`
	Salt     string `json:"-" swaggerignore:"true"`
}

type ForgotPassword struct {
	Email string `json:"email,omitempty"`
}

type ContextInfo struct {
	ID    string
	Role  int
	Email string
}
