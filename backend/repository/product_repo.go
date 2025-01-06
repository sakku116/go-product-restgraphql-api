package repository

import (
	"backend/domain/model"
	"context"
	"errors"
	"fmt"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductRepo struct {
	productColl *qmgo.Collection
}

type IProductRepo interface {
	Create(ctx context.Context, product *model.ProductModel) error
	GetByUUID(ctx context.Context, uuid string) (*model.ProductModel, error)
	Update(ctx context.Context, product *model.ProductModel) error
	Delete(ctx context.Context, uuid string) error
}

func NewProductRepo(productColl *qmgo.Collection) IProductRepo {
	return &ProductRepo{productColl: productColl}
}

func (repo *ProductRepo) Create(ctx context.Context, product *model.ProductModel) error {
	if _, err := repo.productColl.InsertOne(ctx, product); err != nil {
		return fmt.Errorf("failed to create obj")
	}
	return nil
}

func (repo *ProductRepo) GetByUUID(ctx context.Context, uuid string) (*model.ProductModel, error) {
	var product model.ProductModel
	err := repo.productColl.Find(ctx, bson.M{"uuid": uuid}).One(&product)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, errors.New("not found")
		}
		return nil, errors.New("failed to get obj")
	}
	return &product, nil
}

func (repo *ProductRepo) Update(ctx context.Context, product *model.ProductModel) error {
	filter := bson.M{"uuid": product.UUID}
	update := bson.M{"$set": product}

	if err := repo.productColl.UpdateOne(ctx, filter, update); err != nil {
		return errors.New("failed to update user: " + err.Error())
	}
	return nil
}

func (repo *ProductRepo) Delete(ctx context.Context, uuid string) error {
	filter := bson.M{"uuid": uuid}
	if err := repo.productColl.Remove(ctx, filter); err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return errors.New("obj not found")
		}
		return errors.New("failed to delete obj")
	}
	return nil
}
