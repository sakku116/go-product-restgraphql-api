package config

import (
	"context"

	"github.com/qiniu/qmgo"
)

func NewMongoConn(ctx context.Context) *qmgo.Client {
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: Envs.MONGO_URL})
	if err != nil {
		panic(err)
	}
	return client
}
