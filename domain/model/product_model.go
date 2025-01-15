package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductModel struct {
	UUID      string    `json:"uuid" bson:"uuid"`
	UserUUID  string    `json:"user_uuid" bson:"user_uuid"`
	Name      string    `json:"name" bson:"name"`
	Price     float64   `json:"price" bson:"price"`
	Stock     int64     `json:"stock" bson:"stock"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (u ProductModel) GetMongoProps() MongoProps {
	trueTmp := true
	falseTmp := false
	return MongoProps{
		CollName: "products",
		Indexes: []MongoIndex{
			{
				Key: "-uuid",
				Options: &options.IndexOptions{
					Unique:     &trueTmp,
					Background: &trueTmp,
				},
			},
			{
				Key: "-name",
				Options: &options.IndexOptions{
					Unique:     &falseTmp,
					Background: &trueTmp,
				},
			},
		},
	}
}

func (u *ProductModel) ValidateBefore() (err error) {
	// name
	if u.Name == "" {
		return errors.New("name cannot be empty")
	}

	// price
	if u.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	return nil
}
