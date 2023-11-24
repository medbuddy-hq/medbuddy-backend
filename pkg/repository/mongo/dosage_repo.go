package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"time"
)

func (m *Mongo) SaveDosages(ctx context.Context, data []model.Dosage) error {
	db := m.mongoclient.Database(constant.AppName)
	dColl := db.Collection(constant.DosageCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	var records []any
	for i := range data {
		records = append(records, data[i])
	}

	_, err := dColl.InsertMany(ctx, records)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetPatientDosages(ctx context.Context, request *model.DosageFilter) (dosages []model.DosageResponse, err error) {
	db := m.mongoclient.Database(constant.AppName)
	dColl := db.Collection(constant.DosageCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	dosages = []model.DosageResponse{}
	filter := bson.D{
		{Key: "patient_id", Value: request.PatiendID},
	}

	if request.IsActive != nil {
		filter = append(filter, bson.E{
			Key:   "is_active",
			Value: request.IsActive,
		})
	}

	if !request.MedicationID.IsZero() {
		filter = append(filter, bson.E{
			Key:   "medication_id",
			Value: request.MedicationID,
		})
	}

	matchStage := bson.D{{Key: "$match", Value: filter}}
	medicLookupStage, medicUnwindStage := getMedicationLookupAndUnwindStage()
	medLookupStage, medUnwindStage := getDosageMedicineLookupAndUnwindStage()
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"reminder_time", -1}}}}

	pipeline := mongo.Pipeline{matchStage, medicLookupStage, medicUnwindStage, medLookupStage, medUnwindStage, sortStage}
	cur, err := dColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &dosages); err != nil {
		return nil, err
	}

	return dosages, nil
}

func (m *Mongo) SetStatus(ctx context.Context, dosageId, patientId primitive.ObjectID, status string) (found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	dColl := db.Collection(constant.DosageCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	updatesTemp := bson.D{{"is_active", false}, {"status", status}}
	if status == constant.DosageSkipped {
		updatesTemp = append(updatesTemp, bson.E{"time_skipped", time.Now()})
	} else if status == constant.DosageTaken {
		updatesTemp = append(updatesTemp, bson.E{"time_taken", time.Now()})
	}

	filter := bson.D{{Key: "patient_id", Value: patientId}, {Key: "_id", Value: dosageId}, {"is_active", true}}
	update := bson.D{{Key: "$set", Value: updatesTemp}}

	res, err := dColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	if res.MatchedCount <= 0 {
		return false, nil
	}

	return true, nil
}

func (m *Mongo) GetDosage(ctx context.Context, id primitive.ObjectID) (dosage model.DosageResponse, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	dColl := db.Collection(constant.DosageCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}
	medicLookupStage, medicUnwindStage := getMedicationLookupAndUnwindStage()
	medLookupStage, medUnwindStage := getDosageMedicineLookupAndUnwindStage()
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"reminder_time", -1}}}}

	pipeline := mongo.Pipeline{matchStage, medicLookupStage, medicUnwindStage, medLookupStage, medUnwindStage, sortStage}
	cur, err := dColl.Aggregate(ctx, pipeline)
	if err != nil {
		return model.DosageResponse{}, false, err
	}

	var dosages []model.DosageResponse
	if err := cur.All(ctx, &dosages); err != nil {
		return model.DosageResponse{}, false, err
	}

	if len(dosages) < 1 {
		return model.DosageResponse{}, false, nil
	}

	return dosages[0], true, nil
}

func (m *Mongo) DeleteDosages(ctx context.Context, medicationId primitive.ObjectID) (int64, error) {
	db := m.mongoclient.Database(constant.AppName)
	dColl := db.Collection(constant.DosageCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	res, err := dColl.DeleteMany(ctx, bson.D{{Key: "medication_id", Value: medicationId}})
	if err != nil {
		return -1, err
	}

	return res.DeletedCount, nil
}

func getMedicationLookupAndUnwindStage() (medicLookup bson.D, medicUnwind bson.D) {
	medicLookup = bson.D{{
		Key: "$lookup",
		Value: bson.D{{
			Key:   "from",
			Value: "medications",
		}, {
			Key:   "localField",
			Value: "medication_id",
		}, {
			Key:   "foreignField",
			Value: "_id",
		}, {
			Key:   "as",
			Value: "medication",
		}},
	}}

	medicUnwind = bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$medication",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	return
}

func getDosageMedicineLookupAndUnwindStage() (medicineLookup bson.D, medicineUnwind bson.D) {
	medicineLookup = bson.D{{
		Key: "$lookup",
		Value: bson.D{{
			Key:   "from",
			Value: "medicines",
		}, {
			Key:   "localField",
			Value: "medication.medicine_id",
		}, {
			Key:   "foreignField",
			Value: "_id",
		}, {
			Key:   "as",
			Value: "medication.medicine",
		}},
	}}

	medicineUnwind = bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$medication.medicine",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	return
}
