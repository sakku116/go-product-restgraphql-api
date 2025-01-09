package model

import (
	"errors"
	"time"
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

func (u *ProductModel) GetMongoProps() MongoProps {
	return MongoProps{
		CollName: "users",
		Index: []MongoIndex{
			{
				Key:       "uuid",
				Direction: -1,
				Unique:    true,
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
