package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"time"
)

func (m *Mongo) AddTasks(ctx context.Context, tasks []model.Task) (int64, error) {
	db := m.mongoclient.Database(constant.AppName)
	tColl := db.Collection(constant.TaskCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	var records []interface{}
	for i := range tasks {
		records = append(records, tasks[i])
	}

	res, err := tColl.InsertMany(ctx, records)
	if err != nil {
		return -1, err
	}

	return int64(len(res.InsertedIDs)), nil
}

func (m *Mongo) UpdateTask(ctx context.Context, taskID primitive.ObjectID, status string) error {
	db := m.mongoclient.Database(constant.AppName)
	tColl := db.Collection(constant.TaskCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	update := bson.D{{Key: "$set", Value: bson.D{{
		Key:   "status",
		Value: status,
	}}}}
	_, err := tColl.UpdateByID(ctx, taskID, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) DeleteTasks(ctx context.Context, taskIDs []primitive.ObjectID) (int64, error) {
	db := m.mongoclient.Database(constant.AppName)
	tColl := db.Collection(constant.TaskCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: bson.D{{
		Key:   "$in",
		Value: taskIDs,
	}}}}

	res, err := tColl.DeleteMany(ctx, filter)
	if err != nil {
		return -1, err
	}

	return res.DeletedCount, nil
}

func (m *Mongo) GetLatestTasks(ctx context.Context, startTime time.Time) (tasks []model.LatestTaskResponse, err error) {
	db := m.mongoclient.Database(constant.AppName)
	tColl := db.Collection(constant.TaskCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	filter := bson.D{
		{Key: "status", Value: constant.TaskUndone},
		{Key: "time", Value: bson.D{{Key: "$gt", Value: startTime}}},
		{Key: "time", Value: bson.D{{Key: "$lte", Value: startTime.Add(constant.TimeLapseForJobs)}}},
	}

	tasks = []model.LatestTaskResponse{}
	matchStage := bson.D{{Key: "$match", Value: filter}}
	medicLookupStage, medicUnwindStage := getMedicationLookupAndUnwindStage()
	patientLookupStage, patientUnwindStage := getTaskPatientLookupAndUnwindStage()
	medLookupStage, medUnwindStage := getDosageMedicineLookupAndUnwindStage()
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"reminder_time", 1}}}}

	pipeline := mongo.Pipeline{matchStage, medicLookupStage, medicUnwindStage, patientLookupStage, patientUnwindStage,
		medLookupStage, medUnwindStage, sortStage}

	options := options2.Aggregate().SetAllowDiskUse(true)
	cur, err := tColl.Aggregate(ctx, pipeline, options)
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *Mongo) GetTask(ctx context.Context, taskID primitive.ObjectID) (task model.LatestTaskResponse, found bool, err error) {
	db := m.mongoclient.Database(constant.AppName)
	tColl := db.Collection(constant.TaskCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: taskID}}

	var tasks []model.LatestTaskResponse
	matchStage := bson.D{{Key: "$match", Value: filter}}
	medicLookupStage, medicUnwindStage := getMedicationLookupAndUnwindStage()
	patientLookupStage, patientUnwindStage := getTaskPatientLookupAndUnwindStage()
	medLookupStage, medUnwindStage := getDosageMedicineLookupAndUnwindStage()
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"reminder_time", 1}}}}

	pipeline := mongo.Pipeline{matchStage, medicLookupStage, medicUnwindStage, patientLookupStage, patientUnwindStage,
		medLookupStage, medUnwindStage, sortStage}

	cur, err := tColl.Aggregate(ctx, pipeline)
	if err != nil {
		return model.LatestTaskResponse{}, false, err
	}

	err = cur.All(ctx, &tasks)
	if err != nil {
		return model.LatestTaskResponse{}, false, err
	}

	if len(tasks) <= 0 {
		return model.LatestTaskResponse{}, false, nil
	}

	return tasks[0], true, nil
}

func getTaskPatientLookupAndUnwindStage() (patientLookup bson.D, patientUnwind bson.D) {
	patientLookup = bson.D{{
		Key: "$lookup",
		Value: bson.D{{
			Key:   "from",
			Value: "patients",
		}, {
			Key:   "localField",
			Value: "medication.patient_id",
		}, {
			Key:   "foreignField",
			Value: "_id",
		}, {
			Key:   "as",
			Value: "medication.patient",
		}},
	}}

	patientUnwind = bson.D{{
		Key: "$unwind",
		Value: bson.D{{
			Key:   "path",
			Value: "$medication.patient",
		}, {
			Key:   "preserveNullAndEmptyArrays",
			Value: true,
		}},
	}}

	return
}
