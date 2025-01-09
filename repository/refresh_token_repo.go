package repository

import (
	"backend/domain/model"
	"context"
	"errors"
	"fmt"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type RefreshTokenRepo struct {
	refreshTokenColl *qmgo.Collection
}

type IRefreshTokenRepo interface {
	Create(ctx context.Context, refreshToken *model.RefreshTokenModel) error
	GetByToken(ctx context.Context, token string) (*model.RefreshTokenModel, error)
	Update(ctx context.Context, refreshToken *model.RefreshTokenModel) error
	Delete(ctx context.Context, uuid string) error
	InvalidateManyByUserUUID(ctx context.Context, userUUID string) error
}

func NewRefreshTokenRepo(refreshTokenColl *qmgo.Collection) IRefreshTokenRepo {
	return &RefreshTokenRepo{refreshTokenColl: refreshTokenColl}
}

func (repo *RefreshTokenRepo) Create(ctx context.Context, refreshToken *model.RefreshTokenModel) error {
	if _, err := repo.refreshTokenColl.InsertOne(ctx, refreshToken); err != nil {
		return fmt.Errorf("failed to create obj")
	}
	return nil
}

func (repo *RefreshTokenRepo) GetByToken(ctx context.Context, token string) (*model.RefreshTokenModel, error) {
	var refreshToken model.RefreshTokenModel
	err := repo.refreshTokenColl.Find(ctx, bson.M{"token": token}).One(&refreshToken)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get obj")
	}
	return &refreshToken, nil
}

func (repo *RefreshTokenRepo) Update(ctx context.Context, refreshToken *model.RefreshTokenModel) error {
	filter := bson.M{"uuid": refreshToken.UUID}
	update := bson.M{"$set": refreshToken}

	if err := repo.refreshTokenColl.UpdateOne(ctx, filter, update); err != nil {
		return errors.New("failed to update user: " + err.Error())
	}
	return nil
}

func (repo *RefreshTokenRepo) Delete(ctx context.Context, uuid string) error {
	filter := bson.M{"uuid": uuid}
	if err := repo.refreshTokenColl.Remove(ctx, filter); err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return errors.New("obj not found")
		}
		return errors.New("failed to delete obj")
	}
	return nil
}

func (repo *RefreshTokenRepo) InvalidateManyByUserUUID(ctx context.Context, userUUID string) error {
	filter := bson.M{"user_uuid": userUUID}
	update := bson.M{"$set": bson.M{"invalid": true}}
	_, err := repo.refreshTokenColl.UpdateAll(ctx, filter, update)
	if err != nil {
		return errors.New("failed to invalidate refresh tokens: " + err.Error())
	}
	return nil
}
