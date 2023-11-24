package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
)

func (m *Mongo) CreatePractitioner(ctx context.Context, data *model.Practitioner) error {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.PractitionersCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	if _, err := pColl.InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetPractitionerByID(ctx context.Context, id primitive.ObjectID) (pract model.PractitionerResponse, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.PractitionersCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}
	userLookupStage, userUnwindStage := getUserLookupAndUnwindStage()

	pipeline := mongo.Pipeline{matchStage, userLookupStage, userUnwindStage}
	cur, err := pColl.Aggregate(ctx, pipeline)
	if err != nil {
		return model.PractitionerResponse{}, false, err
	}

	var practs []model.PractitionerResponse
	if err := cur.All(ctx, &practs); err != nil {
		return model.PractitionerResponse{}, false, err
	}

	if len(practs) == 0 {
		return model.PractitionerResponse{}, false, nil
	}

	return practs[0], true, nil
}

func (m *Mongo) GetPractitionerByEmail(ctx context.Context, email string) (pract model.PractitionerResponse, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.PractitionersCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "email", Value: email}}}}
	userLookupStage, userUnwindStage := getUserLookupAndUnwindStage()

	pipeline := mongo.Pipeline{matchStage, userLookupStage, userUnwindStage}
	cur, err := pColl.Aggregate(ctx, pipeline)
	if err != nil {
		return model.PractitionerResponse{}, false, err
	}

	var practs []model.PractitionerResponse
	if err := cur.All(ctx, &practs); err != nil {
		return model.PractitionerResponse{}, false, err
	}

	if len(practs) == 0 {
		return model.PractitionerResponse{}, false, nil
	}

	return practs[0], true, nil
}

func (m *Mongo) GetPractitionersByEmail(ctx context.Context, emails []string) (practs []model.PractitionerResponse, err error) {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.PractitionersCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{
		Key: "email", Value: bson.D{{
			Key: "$in", Value: emails,
		}},
	}}}}
	userLookupStage, userUnwindStage := getUserLookupAndUnwindStage()

	pipeline := mongo.Pipeline{matchStage, userLookupStage, userUnwindStage}
	cur, err := pColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	practs = []model.PractitionerResponse{}
	if err := cur.All(ctx, &practs); err != nil {
		return nil, err
	}

	return practs, nil
}

func (m *Mongo) GetPractitionersByIds(ctx context.Context, ids []primitive.ObjectID) (practs []model.PractitionerResponse, err error) {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.PractitionersCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{
		Key: "_id", Value: bson.D{{
			Key: "$in", Value: ids,
		}},
	}}}}
	userLookupStage, userUnwindStage := getUserLookupAndUnwindStage()

	pipeline := mongo.Pipeline{matchStage, userLookupStage, userUnwindStage}
	cur, err := pColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	practs = []model.PractitionerResponse{}
	if err := cur.All(ctx, &practs); err != nil {
		return nil, err
	}

	return practs, nil
}

func (m *Mongo) GetPractitionerMedications(ctx context.Context, practitionerId primitive.ObjectID) (medics []model.MedicationResponse, err error) {
	db := m.mongoclient.Database(constant.AppName)
	pColl := db.Collection(constant.MedicationCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	medics = []model.MedicationResponse{}
	matchStage := bson.D{{Key: "$match", Value: bson.D{{
		Key: "practitioner_ids",
		Value: bson.D{{
			Key: "$elemMatch",
			Value: bson.D{{
				Key:   "$eq",
				Value: practitionerId,
			}},
		}},
	}}}}
	medLookupStage, medUnwindStage := getMedicineLookupAndUnwindStage()
	patientLookupStage, patientUnwindStage := getPatientLookupAndUnwindStage()

	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"created_at", -1}}}}

	pipeline := mongo.Pipeline{matchStage, medLookupStage, medUnwindStage, patientLookupStage, patientUnwindStage, sortStage}
	cur, err := pColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &medics); err != nil {
		return nil, err
	}

	return medics, nil
}

func getPatientLookupAndUnwindStage() (patientLookup bson.D, patientUnwind bson.D) {
	patientLookup = bson.D{{
		Key: "$lookup",
		Value: bson.D{{
			Key:   "from",
			Value: "patients",
		}, {
			Key:   "localField",
			Value: "patient_id",
		}, {
			Key:   "foreignField",
			Value: "_id",
		}, {
			Key:   "as",
			Value: "patient",
		}},
	}}

	patientUnwind = bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$patient",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	return
}
