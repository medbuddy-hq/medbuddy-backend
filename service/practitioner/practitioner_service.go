package practitioner

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/errors"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/pkg/middleware"
	"medbuddy-backend/pkg/repository/storage"
	"medbuddy-backend/utility"
)

type PractitionerService interface {
	CreatePractitioner(data *model.PractitionerRequest) (model.PractitionerResponse, errors.InternalError)
	LoginPractitioner(data *model.UserLogin) (model.PractitionerResponse, errors.InternalError)
	GetPractitioner(uInfo *model.ContextInfo) (model.PractitionerResponse, errors.InternalError)
	GetPractitionerByEmail(email string) (model.PractitionerResponse, errors.InternalError)
	GetPractitionersByIDs(uInfo *model.ContextInfo, ids []string) ([]model.PractitionerResponse, errors.InternalError)
	GetPractitionerMedications(uInfo *model.ContextInfo) ([]model.MedicationResponse, errors.InternalError)
}

type practitionerService struct {
	dbRepo storage.StorageRepository
}

func NewPractitionerService(dbRepo storage.StorageRepository) PractitionerService {
	return &practitionerService{dbRepo: dbRepo}
}

var (
	logger = utility.NewLogger()
)

func (p *practitionerService) CreatePractitioner(data *model.PractitionerRequest) (model.PractitionerResponse, errors.InternalError) {
	ctx := context.Background()

	formatedTime, err := utility.FormatTime(data.DOB)
	if err != nil {
		return model.PractitionerResponse{}, errors.BadRequestError(err.Error())
	}

	hashedPassword, salt, err := utility.HashPassword(data.Password)
	if err != nil {
		logger.Error("Error hashing user's password, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	_, found, err := p.dbRepo.GetPractitionerByEmail(ctx, data.Email)
	if err != nil {
		logger.Error("Error fetching practitioner by email, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	if found {
		return model.PractitionerResponse{}, errors.ResourceNotFoundError("practitioner already exists")
	}

	user := model.User{
		ID:        primitive.NewObjectID(),
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Gender:    data.Gender,
		Email:     data.Email,
		DOB:       formatedTime,
		Role:      constant.Roles[constant.Practitioner],
		Salt:      salt,
		Password:  hashedPassword,
		CreatedAt: utility.ReturnCurrentTime(),
		UpdatedAt: utility.ReturnCurrentTime(),
	}

	if err := p.dbRepo.CreateUser(ctx, &user); err != nil {
		logger.Error("Error creating user document, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	pract := model.Practitioner{
		ID:        primitive.NewObjectID(),
		FullName:  data.Firstname + " " + data.Lastname,
		Title:     data.Title,
		Email:     data.Email,
		UserId:    user.ID,
		Expertise: data.Expertise,
	}

	if err := p.dbRepo.CreatePractitioner(ctx, &pract); err != nil {
		logger.Error("Error creating practitioner's document, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	token, err := middleware.CreateToken(pract.ID.Hex(), pract.Email, user.Role)
	if err != nil {
		logger.Error("Error creating token for user, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	response := utility.RequestsToPractitionerResponse(&pract, &user)
	response.Token = token
	return response, nil
}

func (p *practitionerService) LoginPractitioner(data *model.UserLogin) (model.PractitionerResponse, errors.InternalError) {
	ctx := context.Background()

	practitioner, found, err := p.dbRepo.GetPractitionerByEmail(ctx, data.Email)
	if err != nil {
		logger.Error("Error fetching practitioner by email, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	if !found {
		return model.PractitionerResponse{}, errors.ResourceNotFoundError("practitioner not found")
	}

	if !utility.PasswordIsValid(data.Password, practitioner.User.Salt, practitioner.User.Password) {
		return model.PractitionerResponse{}, errors.BadRequestError("invalid password")
	}

	// Ensure that user is not locked
	if practitioner.User.IsLocked {
		return model.PractitionerResponse{}, errors.ForbiddenError("cannot login, user is currently blocked")
	}

	token, err := middleware.CreateToken(practitioner.ID.Hex(), practitioner.Email, practitioner.User.Role)
	if err != nil {
		logger.Error("Error creating token for user, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	// Omit password and salt from response, then set token
	practitioner.User.Password = ""
	practitioner.User.Salt = ""
	practitioner.Token = token

	return practitioner, nil
}

func (p *practitionerService) GetPractitioner(uInfo *model.ContextInfo) (model.PractitionerResponse, errors.InternalError) {
	ctx := context.Background()

	oId, err := primitive.ObjectIDFromHex(uInfo.ID)
	if err != nil {
		logger.Error("Error converting hex Id to objectId, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	practitioner, found, err := p.dbRepo.GetPractitionerByID(ctx, oId)
	if err != nil {
		logger.Error("Error fetching practitioner by id, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	if !found {
		return model.PractitionerResponse{}, errors.ResourceNotFoundError("practitioner not found")
	}

	return practitioner, nil
}

func (p *practitionerService) GetPractitionerByEmail(email string) (model.PractitionerResponse, errors.InternalError) {
	ctx := context.Background()

	practitioner, found, err := p.dbRepo.GetPractitionerByEmail(ctx, email)
	if err != nil {
		logger.Error("Error fetching practitioner by email, error: ", err.Error())
		return model.PractitionerResponse{}, errors.InternalServerError
	}

	if !found {
		return model.PractitionerResponse{}, errors.ResourceNotFoundError("practitioner not found")
	}

	return practitioner, nil
}

func (p *practitionerService) GetPractitionersByIDs(uInfo *model.ContextInfo, ids []string) ([]model.PractitionerResponse, errors.InternalError) {
	ctx := context.Background()

	var oIds []primitive.ObjectID
	for _, id := range ids {
		oId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			logger.Error("Error converting hex Id to objectId, error: ", err.Error())
			return nil, errors.BadRequestError(fmt.Sprint("invalid id: ", id))
		}

		oIds = append(oIds, oId)
	}

	practitioners, err := p.dbRepo.GetPractitionersByIds(ctx, oIds)
	if err != nil {
		logger.Error("Error fetching practitioners by emails, error: ", err.Error())
		return nil, errors.InternalServerError
	}

	return practitioners, nil
}

func (p *practitionerService) GetPractitionerMedications(userInfo *model.ContextInfo) ([]model.MedicationResponse, errors.InternalError) {
	ctx := context.Background()

	practitionersId, err := primitive.ObjectIDFromHex(userInfo.ID)
	if err != nil {
		logger.Error("Error converting hex Id to objectId, error: ", err.Error())
		return nil, errors.InternalServerError
	}

	medications, err := p.dbRepo.GetPractitionerMedications(ctx, practitionersId)
	if err != nil {
		logger.Error("Error fetching medications for practitioner, error: ", err.Error())
		return nil, errors.InternalServerError
	}

	return medications, nil
}
