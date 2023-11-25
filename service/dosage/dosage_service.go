package dosage

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/errors"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/pkg/repository/storage"
	"medbuddy-backend/utility"
)

type DosageService interface {
	GetPatientsDosages(uInfo model.ContextInfo, isActive *bool, medicationId string) ([]model.DosageResponse, errors.InternalError)
	SetDosageStatus(uInfo *model.ContextInfo, status string, dosageId string) errors.InternalError
	GetDosage(id string) (model.DosageResponse, errors.InternalError)
}

type dosageService struct {
	dbRepo storage.StorageRepository
}

func NewDosageService(dbRepo storage.StorageRepository) DosageService {
	return &dosageService{dbRepo: dbRepo}
}

var (
	logger = utility.NewLogger()
)

func (d *dosageService) GetPatientsDosages(uInfo model.ContextInfo, isActive *bool, medicationId string) ([]model.DosageResponse, errors.InternalError) {
	ctx := context.Background()
	var err error

	oId, err := primitive.ObjectIDFromHex(uInfo.ID)
	if err != nil {
		logger.Error("Error converting hex Id to objectId at GetPatientDosages, error: ", err.Error())
		return nil, errors.InternalServerError
	}

	filter := model.DosageFilter{
		PatiendID: oId,
	}

	var medId primitive.ObjectID
	if medicationId != "" {
		medId, err = primitive.ObjectIDFromHex(medicationId)
		if err != nil {
			logger.Error("Error converting medication id to objectId at GetPatientDosages, error: ", err.Error())
			return nil, errors.BadRequestError("invalid medication id")
		}
		filter.MedicationID = medId
	}

	if isActive != nil {
		filter.IsActive = isActive
	}

	dosages, err := d.dbRepo.GetPatientDosages(ctx, &filter)
	if err != nil {
		logger.Error("Error getting dosages, error: ", err.Error())
		return nil, errors.InternalServerError
	}

	return dosages, nil
}

func (d *dosageService) SetDosageStatus(uInfo *model.ContextInfo, status string, dosageId string) errors.InternalError {
	ctx := context.Background()
	var err error

	oId, err := primitive.ObjectIDFromHex(uInfo.ID)
	if err != nil {
		logger.Error("Error converting hex Id to objectId at GetPatientDosages, error: ", err.Error())
		return errors.InternalServerError
	}

	dId, err := primitive.ObjectIDFromHex(dosageId)
	if err != nil {
		logger.Error("Error converting hex Id to objectId at GetPatientDosages, error: ", err.Error())
		return errors.InternalServerError
	}

	dosage, found, err := d.dbRepo.GetDosage(ctx, dId)
	if err != nil {
		logger.Error("Error fetching dosage document, error: ", err.Error())
		return errors.InternalServerError
	}

	if !found {
		return errors.ResourceNotFoundError("dosage not found")
	}

	if dosage.Status == constant.DosageSkipped || dosage.Status == constant.DosageTaken || !dosage.IsActive {
		return errors.BadRequestError("status of dosage cannot be set again")
	}

	if status != constant.DosageSkipped && status != constant.DosageTaken {
		return errors.BadRequestError("invalid status")
	}

	found, err = d.dbRepo.SetStatus(ctx, dId, oId, status)
	if err != nil {
		logger.Error("Error setting status of dosage, error: ", err.Error())
		return errors.InternalServerError
	}

	if !found {
		return errors.BadRequestError("dosage not found / you don't have access to dosage")
	}

	if status == constant.DosageTaken {
		err := d.dbRepo.IncrementDosageTaken(ctx, dosage.MedicationID)
		if err != nil {
			log.Error("Error when incrementing dosage taken, error: ", err.Error())
		}
	}

	return nil
}

func (d *dosageService) GetDosage(id string) (model.DosageResponse, errors.InternalError) {
	ctx := context.Background()

	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Error converting hex Id to objectId at GetDosage, error: ", err.Error())
		return model.DosageResponse{}, errors.BadRequestError("invalid dosage id")
	}

	dosage, found, err := d.dbRepo.GetDosage(ctx, oId)
	if err != nil {
		logger.Error("Error fetching dosage by id, error: ", err.Error())
		return model.DosageResponse{}, errors.InternalServerError
	}

	if !found {
		return model.DosageResponse{}, errors.ResourceNotFoundError("dosage not found")
	}

	return dosage, nil
}
