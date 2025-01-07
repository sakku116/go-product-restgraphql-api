package repository

import (
	"backend/domain/dto"
	"backend/domain/model"
	"context"
	"errors"
	"fmt"

	"github.com/op/go-logging"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

var logger = logging.MustGetLogger("main")

type ProductRepo struct {
	productColl *qmgo.Collection
}

type IProductRepo interface {
	Create(ctx context.Context, product *model.ProductModel) error
	GetByUUID(ctx context.Context, uuid string) (*model.ProductModel, error)
	GetByName(ctx context.Context, name string, userUUID *string) (*model.ProductModel, error)
	Update(ctx context.Context, product *model.ProductModel) error
	Delete(ctx context.Context, uuid string) error
	GetList(
		ctx context.Context, params dto.ProductRepo_GetListParams,
	) ([]model.ProductModel, *int64, error)
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

func (repo *ProductRepo) GetByName(ctx context.Context, name string, userUUID *string) (*model.ProductModel, error) {
	var product model.ProductModel
	filter := bson.M{"name": name}
	if userUUID != nil {
		filter["user_uuid"] = *userUUID
	}
	err := repo.productColl.Find(ctx, filter).One(&product)
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

func (repo *ProductRepo) GetList(
	ctx context.Context, params dto.ProductRepo_GetListParams,
) ([]model.ProductModel, *int64, error) {
	// validate param
	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}
	var result []model.ProductModel
	var totalCount *int64

	// filter
	matchStage := bson.M{}
	if params.UserUUID != nil {
		matchStage["user_uuid"] = *params.UserUUID
	}
	if params.Query != nil && params.QueryBy != nil {
		matchStage[*params.QueryBy] = bson.M{
			"$regex": *params.Query, "$options": "i",
		}
	}

	// prepare pipelline
	pipeline := []bson.M{}

	// apply match stage to pipeline
	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{"$match": matchStage})
	}

	paginatedResultsPipeline := []bson.M{}

	// sort
	if params.SortBy != nil && *params.SortBy != "" {
		sortStage := bson.M{}
		sortStage[*params.SortBy] = params.SortOrder
		paginatedResultsPipeline = append(paginatedResultsPipeline, bson.M{"$sort": sortStage})
	}

	// pagination
	if params.Page != nil && params.Limit != nil {
		skip := (*params.Page - 1) * *params.Limit
		paginatedResultsPipeline = append(paginatedResultsPipeline, bson.M{"$skip": skip}, bson.M{"$limit": *params.Limit})
	}

	// decide whether to count or not
	if params.DoCount {
		pipeline = append(pipeline, bson.M{
			"$facet": bson.M{
				"paginatedResults": paginatedResultsPipeline,
				"totalCount": []bson.M{
					{"$count": "total_count"},
				},
			},
		})
	} else {
		pipeline = append(pipeline, bson.M{
			"$facet": bson.M{
				"paginatedResults": paginatedResultsPipeline,
			},
		})
	}

	// agggregate
	var aggregateResult []bson.M
	cursor := repo.productColl.Aggregate(ctx, pipeline)
	if err := cursor.All(&aggregateResult); err != nil {
		logger.Error("Error decoding aggregation result:", err)
		return nil, nil, err
	}

	// parse result
	if len(aggregateResult) > 0 {
		data := aggregateResult[0]

		// get paginated result
		if paginatedResultsRaw, exists := data["paginatedResults"]; exists {
			if paginatedResultsArray, ok := paginatedResultsRaw.([]interface{}); ok {
				for _, item := range paginatedResultsArray {
					m, err := bson.Marshal(item)
					if err != nil {
						logger.Error("Error marshalling product:", err)
						return nil, nil, err
					}
					var product model.ProductModel
					if err := bson.Unmarshal(m, &product); err != nil {
						logger.Error("Error unmarshalling product:", err)
						return nil, nil, err
					}
					result = append(result, product)
				}
			}
		}

		// get total count
		if params.DoCount {
			if totalCountRaw, exists := data["totalCount"]; exists {
				if totalCountArray, ok := totalCountRaw.([]interface{}); ok && len(totalCountArray) > 0 {
					if countData, ok := totalCountArray[0].(bson.M); ok {
						if total, exists := countData["total_count"]; exists {
							totalCountVal, ok := total.(int32)
							if ok {
								totalCount = new(int64)
								*totalCount = int64(totalCountVal)
							}
						}
					}
				}
			}
		}
	}

	return result, totalCount, nil
}
