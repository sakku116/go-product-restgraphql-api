package model

import (
	"backend/domain/enum"
	validator_util "backend/utils/validator/user"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
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

func (u UserModel) GetMongoProps() MongoProps {
	trueTmp := true
	falseTmp := false
	return MongoProps{
		CollName: "users",
		Indexes: []MongoIndex{
			{
				Key: "-uuid",
				Options: &options.IndexOptions{
					Unique:     &trueTmp,
					Background: &trueTmp,
				},
			},
			{
				Key: "-username",
				Options: &options.IndexOptions{
					Unique:     &trueTmp,
					Background: &trueTmp,
				},
			},
			{
				Key: "-email",
				Options: &options.IndexOptions{
					Unique:     &falseTmp,
					Background: &trueTmp,
				},
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
