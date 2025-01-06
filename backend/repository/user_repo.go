package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"

	"backend/domain/model"
)

type UserRepo struct {
	userColl *qmgo.Collection
}

type IUserRepo interface {
	Create(ctx context.Context, user *model.UserModel) error
	GetByUUID(ctx context.Context, uuid string) (*model.UserModel, error)
	GetByEmail(ctx context.Context, email string) (*model.UserModel, error)
	GetByUsername(ctx context.Context, username string) (*model.UserModel, error)
	Update(ctx context.Context, user *model.UserModel) error
	Delete(ctx context.Context, uuid string) error
}

func NewUserRepo(userColl *qmgo.Collection) IUserRepo {
	return &UserRepo{userColl: userColl}
}

func (repo *UserRepo) Create(ctx context.Context, user *model.UserModel) error {
	if _, err := repo.userColl.InsertOne(ctx, user); err != nil {
		return fmt.Errorf("failed to create obj")
	}
	return nil
}

func (repo *UserRepo) GetByUUID(ctx context.Context, uuid string) (*model.UserModel, error) {
	var user model.UserModel
	err := repo.userColl.Find(ctx, bson.M{"uuid": uuid}).One(&user)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get obj")
	}
	return &user, nil
}

func (repo *UserRepo) GetByUsername(ctx context.Context, username string) (*model.UserModel, error) {
	var user model.UserModel
	err := repo.userColl.Find(ctx, bson.M{"username": username}).One(&user)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get obj")
	}
	return &user, nil
}

func (repo *UserRepo) GetByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	var user model.UserModel
	err := repo.userColl.Find(ctx, bson.M{"email": email}).One(&user)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get obj")
	}
	return &user, nil
}

func (repo *UserRepo) Update(ctx context.Context, user *model.UserModel) error {
	filter := bson.M{"uuid": user.UUID}
	update := bson.M{"$set": user}

	if err := repo.userColl.UpdateOne(ctx, filter, update); err != nil {
		return errors.New("failed to update obj")
	}
	return nil
}

func (repo *UserRepo) Delete(ctx context.Context, uuid string) error {
	filter := bson.M{"uuid": uuid}
	if err := repo.userColl.Remove(ctx, filter); err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return errors.New("obj not found")
		}
		return errors.New("failed to delete obj")
	}
	return nil
}
