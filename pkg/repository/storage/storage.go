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
}
