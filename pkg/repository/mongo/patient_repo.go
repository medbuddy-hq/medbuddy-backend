package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
)

func (m *Mongo) CreatePatient(ctx context.Context, data *model.Patient) error {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.PatientsCollection)

	_, found, err := m.GetPatientByEmail(ctx, data.Email)
	if err != nil {
		return err
	}

	if found {
		return constant.ErrResourceAlreadyExists
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	if _, err := pColl.InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetPatientByID(ctx context.Context, id primitive.ObjectID) (patient model.PatientResponse, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.PatientsCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}
	userLookupStage, userUnwindStage := getUserLookupAndUnwindStage()

	pipeline := mongo.Pipeline{matchStage, userLookupStage, userUnwindStage}
	cur, err := pColl.Aggregate(ctx, pipeline)
	if err != nil {
		return model.PatientResponse{}, false, err
	}

	var patients []model.PatientResponse
	if err := cur.All(ctx, &patients); err != nil {
		return model.PatientResponse{}, false, err
	}

	if len(patients) == 0 {
		return model.PatientResponse{}, false, nil
	}

	return patients[0], true, nil
}

func (m *Mongo) GetPatientByEmail(ctx context.Context, email string) (patient model.PatientResponse, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.PatientsCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "email", Value: email}}}}
	userLookupStage, userUnwindStage := getUserLookupAndUnwindStage()

	pipeline := mongo.Pipeline{matchStage, userLookupStage, userUnwindStage}
	cur, err := pColl.Aggregate(ctx, pipeline)
	if err != nil {
		return model.PatientResponse{}, false, err
	}

	var patients []model.PatientResponse
	if err := cur.All(ctx, &patients); err != nil {
		return model.PatientResponse{}, false, err
	}

	if len(patients) == 0 {
		return model.PatientResponse{}, false, nil
	}

	return patients[0], true, nil
}

func getUserLookupAndUnwindStage() (userLookup bson.D, userUnwind bson.D) {
	userLookup = bson.D{{
		Key: "$lookup",
		Value: bson.D{{
			Key:   "from",
			Value: "users",
		}, {
			Key:   "localField",
			Value: "user_id",
		}, {
			Key:   "foreignField",
			Value: "_id",
		}, {
			Key:   "as",
			Value: "user",
		}},
	}}

	userUnwind = bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$user",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	return
}
