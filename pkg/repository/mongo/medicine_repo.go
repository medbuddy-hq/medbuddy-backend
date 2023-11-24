package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
)

func (m *Mongo) AddMedicine(ctx context.Context, data *model.Medicine) error {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicineCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	if _, err := mColl.InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) GetMedicineByID(ctx context.Context, id primitive.ObjectID) (medicine model.Medicine, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicineCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	filter := bson.D{{"_id", id}}
	if err := mColl.FindOne(ctx, filter).Decode(&medicine); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Medicine{}, false, nil
		}
		return model.Medicine{}, false, err
	}

	return medicine, true, nil
}

func (m *Mongo) GetMedicineFilter(ctx context.Context, req *model.MedicineFilter) (medicine model.Medicine, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicineCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	filter := bson.D{
		{Key: "name", Value: req.Name},
		{Key: "manufacturer", Value: req.Manufacturer},
		{Key: "strength", Value: req.Strength},
	}

	if req.Form != "" {
		filter = append(filter, bson.E{Key: "form", Value: req.Form})
	}

	if err := mColl.FindOne(ctx, filter).Decode(&medicine); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Medicine{}, false, nil
		}
		return model.Medicine{}, false, err
	}

	return medicine, true, nil
}

func (m *Mongo) UpdateMedicine(ctx context.Context, id primitive.ObjectID, data *model.Medicine) (found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicineCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	update := bson.D{{"$set", data}}
	res, err := mColl.UpdateByID(ctx, id, update)
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

func (m *Mongo) DeleteMedicine(ctx context.Context, id primitive.ObjectID) (found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	mColl := db.Collection(constant.MedicineCollection)

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
