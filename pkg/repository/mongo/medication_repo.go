package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
)

func (m *Mongo) AddMedication(ctx context.Context, data *model.Medication) error {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicationCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	if _, err := mColl.InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) UpdateMedication(ctx context.Context, id primitive.ObjectID, data *model.MedicationRequest) (found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicationCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	res, err := mColl.UpdateByID(ctx, id, data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	if res.MatchedCount < 1 {
		return false, nil
	}

	return true, nil
}

func (m *Mongo) DeleteMedication(ctx context.Context, id primitive.ObjectID) (found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicationCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	filter := bson.D{{"_id", id}}
	res, err := mColl.DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	if res.DeletedCount < 1 {
		return false, nil
	}

	return true, nil
}

func (m *Mongo) GetMedication(ctx context.Context, id primitive.ObjectID) (medic model.MedicationResponse, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicationCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	medics := []model.MedicationResponse{}
	matchStage := bson.D{{Key: "$match", Value: bson.D{{"_id", id}}}}
	medLookupStage, medUnwindStage := getMedicineLookupAndUnwindStage()

	pipeline := mongo.Pipeline{matchStage, medLookupStage, medUnwindStage}
	cur, err := mColl.Aggregate(ctx, pipeline)
	if err != nil {
		return model.MedicationResponse{}, false, err
	}

	if err := cur.All(ctx, &medics); err != nil {
		return model.MedicationResponse{}, false, err
	}

	if len(medics) <= 0 {
		return model.MedicationResponse{}, false, nil
	}

	return medics[0], true, nil
}

func (m *Mongo) GetPatientsMedications(ctx context.Context, patientId primitive.ObjectID) (medics []model.MedicationResponse, err error) {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicationCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	medics = []model.MedicationResponse{}
	matchStage := bson.D{{Key: "$match", Value: bson.D{{"patient_id", patientId}}}}
	medLookupStage, medUnwindStage := getMedicineLookupAndUnwindStage()
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"created_at", -1}}}}

	pipeline := mongo.Pipeline{matchStage, medLookupStage, medUnwindStage, sortStage}
	cur, err := mColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &medics); err != nil {
		return nil, err
	}

	return medics, nil
}

func getMedicineLookupAndUnwindStage() (medicineLookup bson.D, medicineUnwind bson.D) {
	medicineLookup = bson.D{{
		Key: "$lookup",
		Value: bson.D{{
			Key:   "from",
			Value: "medicines",
		}, {
			Key:   "localField",
			Value: "medicine_id",
		}, {
			Key:   "foreignField",
			Value: "_id",
		}, {
			Key:   "as",
			Value: "medicine",
		}},
	}}

	medicineUnwind = bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$medicine",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	return
}
