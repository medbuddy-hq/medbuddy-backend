package patient

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/errors"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/pkg/middleware"
	"medbuddy-backend/pkg/repository/storage"
	"medbuddy-backend/utility"
)

type PatientService interface {
	CreatePatient(data *model.CreatePatientReq) (model.PatientResponse, errors.InternalError)
	LoginPatient(data *model.UserLogin) (model.PatientResponse, errors.InternalError)
	GetPatient(id string) (model.PatientResponse, errors.InternalError)
	GetPatientByEmail(email string) (model.PatientResponse, errors.InternalError)
}

type patientService struct {
	dbRepo storage.StorageRepository
}

func NewPatientService(dbRepo storage.StorageRepository) PatientService {
	return &patientService{dbRepo: dbRepo}
}

var (
	logger = utility.NewLogger()
)

func (p *patientService) CreatePatient(data *model.CreatePatientReq) (model.PatientResponse, errors.InternalError) {
	formatedTime, err := utility.FormatTime(data.DOB)
	if err != nil {
		return model.PatientResponse{}, errors.BadRequestError(err.Error())
	}

	hashedPassword, salt, err := utility.HashPassword(data.Password)
	if err != nil {
		logger.Error("Error hashing user's password, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	user := model.User{
		ID:        primitive.NewObjectID(),
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Gender:    data.Gender,
		Email:     data.Email,
		DOB:       formatedTime,
		Role:      constant.Roles[constant.Patient],
		Salt:      salt,
		Password:  hashedPassword,
		CreatedAt: utility.ReturnCurrentTime(),
		UpdatedAt: utility.ReturnCurrentTime(),
	}

	ctx := context.Background()
	if err := p.dbRepo.CreateUser(ctx, &user); err != nil {
		logger.Error("Error creating user document, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	patient := model.Patient{
		ID:       primitive.NewObjectID(),
		FullName: data.Firstname + " " + data.Lastname,
		Email:    data.Email,
		UserID:   user.ID,
	}

	_, found, err := p.dbRepo.GetPatientByEmail(ctx, data.Email)
	if err != nil {
		logger.Error("Error fetching patient by email, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	if found {
		return model.PatientResponse{}, errors.ResourceNotFoundError("patient already exists")
	}

	if err := p.dbRepo.CreatePatient(ctx, &patient); err != nil {
		logger.Error("Error creating patient's document, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	token, err := middleware.CreateToken(patient.ID.Hex(), patient.Email, user.Role)
	if err != nil {
		logger.Error("Error creating token for user, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	response := utility.RequestsToPatientResponse(&patient, &user)
	response.Token = token
	return response, nil
}

func (p *patientService) LoginPatient(data *model.UserLogin) (model.PatientResponse, errors.InternalError) {
	// Get database details
	ctx := context.Background()

	patient, found, err := p.dbRepo.GetPatientByEmail(ctx, data.Email)
	if err != nil {
		logger.Error("Error fetching patient by email, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	if !found {
		return model.PatientResponse{}, errors.ResourceNotFoundError("patient not found")
	}

	if !utility.PasswordIsValid(data.Password, patient.User.Salt, patient.User.Password) {
		return model.PatientResponse{}, errors.BadRequestError("invalid password")
	}

	// Ensure that user is not locked
	if patient.User.IsLocked {
		return model.PatientResponse{}, errors.ForbiddenError("cannot login, user is currently blocked")
	}

	token, err := middleware.CreateToken(patient.ID.Hex(), patient.Email, patient.User.Role)
	if err != nil {
		logger.Error("Error creating token for user, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	// Omit password and salt from response, then set token
	patient.User.Password = ""
	patient.User.Salt = ""
	patient.Token = token

	return patient, nil
}

func (p *patientService) GetPatient(id string) (model.PatientResponse, errors.InternalError) {
	ctx := context.Background()

	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Error converting hex Id to objectId, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	patient, found, err := p.dbRepo.GetPatientByID(ctx, oId)
	if err != nil {
		logger.Error("Error fetching patient by id, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	if !found {
		return model.PatientResponse{}, errors.ResourceNotFoundError("patient not found")
	}

	return patient, nil
}

func (p *patientService) GetPatientByEmail(email string) (model.PatientResponse, errors.InternalError) {
	ctx := context.Background()

	patient, found, err := p.dbRepo.GetPatientByEmail(ctx, email)
	if err != nil {
		logger.Error("Error fetching patient by email, error: ", err.Error())
		return model.PatientResponse{}, errors.InternalServerError
	}

	if !found {
		return model.PatientResponse{}, errors.ResourceNotFoundError("patient not found")
	}

	return patient, nil
}
