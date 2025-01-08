package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"

	"backend/domain/dto"
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
	GetList(
		ctx context.Context, params dto.UserRepo_GetListParams,
	) ([]model.UserModel, int64, error)
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

func (repo *UserRepo) GetList(
	ctx context.Context, params dto.UserRepo_GetListParams,
) ([]model.UserModel, int64, error) {
	// validate param
	err := params.Validate()
	if err != nil {
		return nil, 0, err
	}
	var result []model.UserModel
	var totalCount int64

	// filter
	matchStage := bson.M{}
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

	// logger.Debugf("pipeline: %v", helper.PrettyJson(pipeline))

	// agggregate
	var aggregateResult []bson.M
	cursor := repo.userColl.Aggregate(ctx, pipeline)
	if err := cursor.All(&aggregateResult); err != nil {
		logger.Error("Error decoding aggregation result:", err)
		return nil, 0, err
	}

	// logger.Debugf("aggregateResult: %v", helper.PrettyJson(aggregateResult))

	// parse result
	if len(aggregateResult) > 0 {
		data := aggregateResult[0]

		// get paginated result
		if paginatedResultsRaw, exists := data["paginatedResults"]; exists {
			if paginatedResultsArray, ok := paginatedResultsRaw.(bson.A); ok {
				for _, item := range paginatedResultsArray {
					logger.Debug("item")
					m, err := bson.Marshal(item)
					if err != nil {
						logger.Error("Error marshalling user:", err)
						return nil, 0, err
					}
					var user model.UserModel
					if err := bson.Unmarshal(m, &user); err != nil {
						logger.Error("Error unmarshalling user:", err)
						return nil, 0, err
					}
					result = append(result, user)
				}
			}
		}

		// get total count
		if params.DoCount {
			if totalCountRaw, exists := data["totalCount"]; exists {
				if totalCountArray, ok := totalCountRaw.(bson.A); ok && len(totalCountArray) > 0 {
					if countData, ok := totalCountArray[0].(bson.M); ok {
						if total, exists := countData["total_count"]; exists {
							totalCountVal, ok := total.(int32)
							if ok {
								totalCount = int64(totalCountVal)
							}
						}
					}
				}
			}
		}
	}

	return result, totalCount, nil
}
