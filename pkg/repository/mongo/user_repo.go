package mongo

import (
	"context"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
)

func (m *Mongo) CreateUser(ctx context.Context, data *model.User) error {
	db := m.mongoclient.Database(constant.AppName)
	uColl := db.Collection(constant.UsersCollection)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, m.timeout)
	defer cancel()

	if _, err := uColl.InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}
