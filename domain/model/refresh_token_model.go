package model

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshTokenModel struct {
	UUID      string     `json:"uuid" bson:"uuid"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	UserUUID  string     `json:"user_uuid" bson:"user_uuid"`
	Token     string     `json:"token" bson:"token"`
	UsedAt    *time.Time `json:"used_at" bson:"used_at"`
	ExpiredAt *time.Time `json:"expired_at" bson:"expired_at"`
	Invalid   bool       `json:"invalid" bson:"invalid"`
}

func (u RefreshTokenModel) GetMongoProps() MongoProps {
	trueTmp := true
	falseTmp := false
	return MongoProps{
		CollName: "refresh_tokens",
		Indexes: []MongoIndex{
			{
				Key: "-uuid",
				Options: &options.IndexOptions{
					Unique:     &trueTmp,
					Background: &trueTmp,
				},
			},
			{
				Key: "-token",
				Options: &options.IndexOptions{
					Unique:     &trueTmp,
					Background: &trueTmp,
				},
			},
			{
				Key: "-user_uuid",
				Options: &options.IndexOptions{
					Unique:     &falseTmp,
					Background: &trueTmp,
				},
			},
		},
	}
}
