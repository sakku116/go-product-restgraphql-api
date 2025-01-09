package model

import (
	"time"
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

func (u *RefreshTokenModel) GetMongoProps() MongoProps {
	return MongoProps{
		CollName: "refresh_tokens",
		Index: []MongoIndex{
			{
				Key:       "uuid",
				Direction: -1,
				Unique:    true,
			},
			{
				Key:       "token",
				Direction: -1,
				Unique:    true,
			},
		},
	}
}
