package model

import (
	"backend/domain/enum"
	validator_util "backend/utils/validator/user"
	"errors"
	"time"
)

type UserModel struct {
	UUID      string        `json:"uuid" bson:"uuid"`
	Username  string        `json:"username" bson:"username"`
	Password  string        `json:"password" bson:"password"`
	Role      enum.UserRole `json:"role" bson:"role"`
	Email     string        `json:"email" bson:"email"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

func (u *UserModel) GetMongoProps() MongoProps {
	return MongoProps{
		CollName: "users",
		Index: []MongoIndex{
			{
				Key:       "username",
				Direction: -1,
				Unique:    true,
			},
			{
				Key:       "uuid",
				Direction: -1,
				Unique:    true,
			},
		},
	}
}

func (u *UserModel) ValidateBefore() (err error) {
	// username
	err = validator_util.ValidateUsername(u.Username)
	if err != nil {
		return errors.New("user validation error: " + err.Error())
	}

	// email
	err = validator_util.ValidateEmail(u.Email)
	if err != nil {
		return errors.New("user validation error: " + err.Error())
	}

	// password
	err = validator_util.ValidatePassword(u.Password)
	if err != nil {
		return errors.New("user validation error: " + err.Error())
	}

	return nil
}
