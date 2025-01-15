package model

import "go.mongodb.org/mongo-driver/mongo/options"

type MongoIndex struct {
	Key     string
	Options *options.IndexOptions
}

type MongoProps struct {
	CollName string
	Indexes  []MongoIndex
}

type IModel interface {
	GetMongoProps() MongoProps
}
