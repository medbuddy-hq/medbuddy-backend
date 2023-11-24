package medicine

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medbuddy-backend/internal/errors"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/pkg/repository/storage"
	"medbuddy-backend/utility"
)

type MedicineService interface {
	AddMedicine(data *model.MedicineRequest) (model.Medicine, errors.InternalError)
	GetMedicine(id string) (model.Medicine, errors.InternalError)
	GetMedicineFilter(req *model.MedicineFilter) (model.Medicine, errors.InternalError)
	UpdateMedicine(id string, data *model.MedicineRequest) (model.Medicine, errors.InternalError)
	DeleteMedicine(id string) errors.InternalError
}

type medicineService struct {
	dbRepo storage.StorageRepository
}

func NewMedicineService(dbRepo storage.StorageRepository) MedicineService {
	return &medicineService{dbRepo: dbRepo}
}

var (
	logger = utility.NewLogger()
)

func (m *medicineService) AddMedicine(data *model.MedicineRequest) (model.Medicine, errors.InternalError) {
	data.ID = primitive.NewObjectID()
	data.CreatedAt = utility.ReturnCurrentTime()
	data.UpdatedAt = utility.ReturnCurrentTime()
	medicine := utility.MedicineRequestToMedicine(data)
	ctx := context.Background()

	medFilter := model.MedicineFilter{
		Name:         data.Name,
		Manufacturer: data.Manufacturer,
		Strength:     data.Strength,
		Form:         data.Form,
	}
	_, found, err := m.dbRepo.GetMedicineFilter(ctx, &medFilter)
	if err != nil {
		logger.Error("Error fetching medicine by filters, error: ", err.Error())
		return model.Medicine{}, errors.InternalServerError
	}

	if found {
		return model.Medicine{}, errors.BadRequestError("medicine already exists")
	}

	if err := m.dbRepo.AddMedicine(ctx, &medicine); err != nil {
		logger.Error("Error adding medicine, error: ", err.Error())
		return model.Medicine{}, errors.InternalServerError
	}

	return medicine, nil
}

func (m *medicineService) GetMedicine(id string) (model.Medicine, errors.InternalError) {
	ctx := context.Background()

	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Error converting hex Id to objectId, error: ", err.Error())
		return model.Medicine{}, errors.BadRequestError("invalid medicine id")
	}

	logger.Info("oid: ", oId)

	medicine, found, err := m.dbRepo.GetMedicineByID(ctx, oId)
	if err != nil {
		logger.Error("Error fetching medicine by id, error: ", err.Error())
		return model.Medicine{}, errors.InternalServerError
	}

	if !found {
		return model.Medicine{}, errors.ResourceNotFoundError("medicine not found")
	}

	return medicine, nil
}

func (m *medicineService) GetMedicineFilter(req *model.MedicineFilter) (model.Medicine, errors.InternalError) {
	ctx := context.Background()

	medicine, found, err := m.dbRepo.GetMedicineFilter(ctx, req)
	if err != nil {
		logger.Error("Error fetching medicine by filters, error: ", err.Error())
		return model.Medicine{}, errors.InternalServerError
	}

	if !found {
		return model.Medicine{}, errors.ResourceNotFoundError("medicine not found")
	}

	return medicine, nil
}

func (m *medicineService) UpdateMedicine(id string, data *model.MedicineRequest) (model.Medicine, errors.InternalError) {
	ctx := context.Background()

	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Error converting hex Id to objectId, error: ", err.Error())
		return model.Medicine{}, errors.InternalServerError
	}

	data.ID = oId
	data.UpdatedAt = utility.ReturnCurrentTime()

	medicine := utility.MedicineRequestToMedicine(data)
	found, err := m.dbRepo.UpdateMedicine(ctx, oId, &medicine)
	if err != nil {
		logger.Error("Error updating medicine by id, error: ", err.Error())
		return model.Medicine{}, errors.InternalServerError
	}

	if !found {
		return model.Medicine{}, errors.ResourceNotFoundError("medicine not found")
	}

	return medicine, nil
}

func (m *medicineService) DeleteMedicine(id string) errors.InternalError {
	ctx := context.Background()

	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Error converting hex Id to objectId, error: ", err.Error())
		return errors.InternalServerError
	}

	found, err := m.dbRepo.DeleteMedicine(ctx, oId)
	if err != nil {
		logger.Error("Error deleting medicine by id, error: ", err.Error())
		return errors.InternalServerError
	}

	if !found {
		return errors.ResourceNotFoundError("medicine not found")
	}

	return nil
}
