package mongo_util

import (
	"backend/domain/model"
	"context"

	"github.com/op/go-logging"
	"github.com/qiniu/qmgo"
	qmgo_opts "github.com/qiniu/qmgo/options"
)

var logger = logging.MustGetLogger("mongo_util")

func EnsureMongoIndexes(database *qmgo.Database, models ...model.IModel) error {
	for _, model := range models {
		collName := model.GetMongoProps().CollName
		logger.Infof("ensure index for %s", collName)
		if collName == "" {
			logger.Warningf("collection name empty")
		}
		collection := database.Collection(collName)
		for _, index := range model.GetMongoProps().Indexes {
			if err := collection.CreateOneIndex(context.Background(), qmgo_opts.IndexModel{
				Key:          []string{index.Key},
				IndexOptions: index.Options,
			}); err != nil {
				logger.Errorf("ensure index error: %v", err)
			}
		}
		logger.Infof("ensure index for %s done", collName)
	}
	return nil
}
