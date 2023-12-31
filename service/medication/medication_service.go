package medication

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/errors"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/pkg/repository/storage"
	"medbuddy-backend/utility"
)

type MedicationService interface {
	AddMedication(userInfo *model.ContextInfo, data *model.MedicationRequest) (model.MedicationResponse, errors.InternalError)
	GetMedication(id string) (model.MedicationResponse, errors.InternalError)
	GetPatientMedications(userInfo model.ContextInfo) ([]model.MedicationResponse, errors.InternalError)
	UpdateMedication(userInfo *model.ContextInfo, id string, data *model.MedicationRequest) (model.MedicationResponse, errors.InternalError)
	DeleteMedication(userInfo *model.ContextInfo, id string) errors.InternalError
	AddPractitionersToMedication(userInfo *model.ContextInfo, medicId string, practEmails []string) (string, errors.InternalError)
}

type medicationService struct {
	dbRepo storage.StorageRepository
}

func NewMedicationService(dbRepo storage.StorageRepository) MedicationService {
	return &medicationService{dbRepo: dbRepo}
}

var (
	logger = utility.NewLogger()
)

func (m *medicationService) AddMedication(userInfo *model.ContextInfo, data *model.MedicationRequest) (model.MedicationResponse, errors.InternalError) {
	ctx := context.Background()

	patientID, err := primitive.ObjectIDFromHex(userInfo.ID)
	if err != nil {
		logger.Error("Error converting hex Id to objectId at AddMedication, error: ", err.Error())
		return model.MedicationResponse{}, errors.InternalServerError
	}

	startDate, err := utility.FormatTime(data.StartDate)
	if err != nil {
		logger.Error("Error converting startDate in AddMedication, error: ", err.Error())
		return model.MedicationResponse{}, errors.BadRequestError(fmt.Sprint("StartDate: ", err.Error()))
	}

	medication := model.Medication{
		ID:                  primitive.NewObjectID(),
		Name:                data.Name,
		PatientID:           patientID,
		StartDate:           startDate,
		DosageQuantity:      data.DosageQuantity,
		DailyDosage:         data.DailyDosage,
		Treatment:           data.Treatment,
		CreatedAt:           utility.ReturnCurrentTime(),
		UpdatedAt:           utility.ReturnCurrentTime(),
		IsActive:            true,
		Comment:             data.Comment,
		TotalNumberOfDosage: data.TotalNumberOfDosage,
	}

	if medication.TotalNumberOfDosage <= 0 {
		return model.MedicationResponse{}, errors.BadRequestError("invalid value for total number of dosage")
	}

	med, found, err := m.dbRepo.GetMedicineFilter(ctx, &model.MedicineFilter{
		Name:         data.Medicine.Name,
		Manufacturer: data.Medicine.Manufacturer,
		Strength:     data.Medicine.Strength,
		Form:         data.Medicine.Form,
	})

	if err != nil {
		logger.Error("Error checking if medicine exists in AddMedication, error: ", err.Error())
		return model.MedicationResponse{}, errors.InternalServerError
	}

	if found {
		medication.MedicineID = med.ID
	} else {
		data.Medicine.ID = primitive.NewObjectID()
		data.Medicine.CreatedAt = utility.ReturnCurrentTime()
		data.Medicine.UpdatedAt = utility.ReturnCurrentTime()

		if err := m.dbRepo.AddMedicine(ctx, &data.Medicine); err != nil {
			logger.Error("Error adding new medicine in AddMedication, error: ", err.Error())
			return model.MedicationResponse{}, errors.InternalServerError
		}

		medication.MedicineID = data.Medicine.ID
	}

	dosages, err := utility.GetDosages(medication.StartDate, data)
	if err != nil {
		logger.Error("Error getting dosage times in AddMedication, error: ", err.Error())
		return model.MedicationResponse{}, errors.BadRequestError(err.Error())
	}

	for i := range dosages {
		dosages[i].ID = primitive.NewObjectID()
		dosages[i].MedicationID = medication.ID
		dosages[i].PatientID = patientID
	}

	if err := m.dbRepo.SaveDosages(ctx, dosages); err != nil {
		logger.Error("Error saving dosages in AddMedication, error: ", err.Error())
		return model.MedicationResponse{}, errors.InternalServerError
	}

	if err := m.dbRepo.AddMedication(ctx, &medication); err != nil {
		logger.Error("Error adding medication in AddMedication, error: ", err.Error())
		return model.MedicationResponse{}, errors.InternalServerError
	}

	var tasks []model.Task
	for _, dosage := range dosages {
		t := model.Task{
			ID:           primitive.NewObjectID(),
			Time:         dosage.ReminderTime,
			Status:       constant.TaskUndone,
			MedicationID: medication.ID,
		}

		tasks = append(tasks, t)
	}

	count, err := m.dbRepo.AddTasks(ctx, tasks)
	if err != nil {
		logger.Error("Error adding tasks in AddMedication, error: ", err.Error())
		return model.MedicationResponse{}, errors.InternalServerError
	}

	logger.Infof("Successfully added '%v' task(s) for %v medication", count, data.Medicine.Name)

	response := utility.MedicationToMedicationResponse(&medication)
	response.Dosages = dosages
	response.Medicine = data.Medicine
	response.Medicine.ID = medication.MedicineID
	response.Patient = model.Patient{ID: patientID, Email: userInfo.Email}

	return response, nil
}

func (m *medicationService) GetMedication(id string) (model.MedicationResponse, errors.InternalError) {
	ctx := context.Background()

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Error converting hex Id to objectId at GetMedication, error: ", err.Error())
		return model.MedicationResponse{}, errors.BadRequestError("invalid id")
	}

	medic, found, err := m.dbRepo.GetMedication(ctx, oID)
	if err != nil {
		logger.Error("Error fetching medication by id, error: ", err.Error())
		return model.MedicationResponse{}, errors.InternalServerError
	}

	if !found {
		return model.MedicationResponse{}, errors.ResourceNotFoundError("medication does not exist")
	}

	return medic, nil
}

func (m *medicationService) GetPatientMedications(userInfo model.ContextInfo) ([]model.MedicationResponse, errors.InternalError) {
	ctx := context.Background()

	oId, err := primitive.ObjectIDFromHex(userInfo.ID)
	if err != nil {
		logger.Error("Error converting hex Id to objectId at GetPatientMedications error: ", err.Error())
		return nil, errors.InternalServerError
	}

	medics, err := m.dbRepo.GetPatientsMedications(ctx, oId)
	if err != nil {
		logger.Error("Error getting patients medications, error: ", err.Error())
		return nil, errors.InternalServerError
	}

	return medics, nil
}

func (m *medicationService) UpdateMedication(userInfo *model.ContextInfo, id string, data *model.MedicationRequest) (model.MedicationResponse, errors.InternalError) {
	//ctx := context.Background()

	return model.MedicationResponse{}, nil
}

func (m *medicationService) DeleteMedication(userInfo *model.ContextInfo, id string) errors.InternalError {
	ctx := context.Background()

	medId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Error converting hex Id to objectId, error: ", err.Error())
		return errors.InternalServerError
	}

	count, err := m.dbRepo.DeleteDosages(ctx, medId)
	if err != nil {
		logger.Error("Error deleting dosages, error: ", err.Error())
		return errors.InternalServerError
	}

	logger.Infof("Matched and deleted %v dosage(s)", count)

	found, err := m.dbRepo.DeleteMedication(ctx, medId)
	if err != nil {
		logger.Error("Error deleting medication by id, error: ", err.Error())
		return errors.InternalServerError
	}

	if !found {
		return errors.ResourceNotFoundError("medication not found")
	}

	return nil
}

func (m *medicationService) AddPractitionersToMedication(userInfo *model.ContextInfo, medicId string, practEmails []string) (string, errors.InternalError) {
	ctx := context.Background()

	medId, err := primitive.ObjectIDFromHex(medicId)
	if err != nil {
		logger.Error("Error converting hex Id to objectId, error: ", err.Error())
		return "", errors.InternalServerError
	}

	medication, found, err := m.dbRepo.GetMedication(ctx, medId)
	if err != nil {
		logger.Error("Error fetching medication, error: ", err.Error())
		return "", errors.InternalServerError
	}

	if !found {
		return "", errors.ResourceNotFoundError("medication not found")
	}

	practitioners, err := m.dbRepo.GetPractitionersByEmail(ctx, practEmails)
	if err != nil {
		logger.Error("Error fetching specified practioner(s), error: ", err.Error())
		return "", errors.InternalServerError
	}

	if len(practitioners) <= 0 {
		return "", errors.BadRequestError("invalid practitioner email(s)")
	}
	logger.Infof("Successfully fetched %v out of %v practitioners\n", len(practitioners), len(practEmails))

	practitionerIds := make([]primitive.ObjectID, len(practitioners))
	for i, p := range practitioners {
		practitionerIds[i] = p.ID
	}
	medication.PractitionerIDs = append(medication.PractitionerIDs, practitionerIds...)

	if _, err := m.dbRepo.AddPractitionerToMed(ctx, medId, medication.PractitionerIDs); err != nil {
		logger.Error("Error adding practioner(s) to medication, error: ", err.Error())
		return "", errors.InternalServerError
	}

	return fmt.Sprintf("successfully added %v out of %v practitioner(s) to medication", len(practitionerIds), len(practEmails)), nil
}
