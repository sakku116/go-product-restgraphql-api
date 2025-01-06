package model

type MongoIndex struct {
	Key       string
	Unique    bool
	Direction int
}

type MongoProps struct {
	CollName string
	Index    []MongoIndex
}
