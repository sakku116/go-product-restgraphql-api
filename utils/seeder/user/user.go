package seeder_util

import (
	"backend/config"
	"backend/domain/enum"
	"backend/domain/model"
	"backend/repository"
	bcrypt_util "backend/utils/bcrypt"
	"backend/utils/helper"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

func SeedUser(userRepo repository.IUserRepo) error {
	ctx := context.Background()
	users := []model.UserModel{}

	if config.Envs.INITIAL_ADMIN_USERNAME != "" && config.Envs.INITIAL_ADMIN_PASSWORD != "" {
		hashedPassword, _ := bcrypt_util.Hash(config.Envs.INITIAL_ADMIN_PASSWORD)
		users = append(users, model.UserModel{
			UUID:      uuid.New().String(),
			Username:  config.Envs.INITIAL_ADMIN_USERNAME,
			Password:  hashedPassword,
			Email:     fmt.Sprint(config.Envs.INITIAL_ADMIN_USERNAME, "@gmail.com"),
			Role:      enum.UserRole_Admin,
			CreatedAt: helper.TimeNowUTC(),
			UpdatedAt: helper.TimeNowUTC(),
		})
	} else {
		logger.Warningf("initial admin username and password not set")
	}

	if config.Envs.INITIAL_USER_USERNAME != "" && config.Envs.INITIAL_USER_PASSWORD != "" {
		hashedPassword, _ := bcrypt_util.Hash(config.Envs.INITIAL_USER_PASSWORD)
		users = append(users, model.UserModel{
			UUID:      uuid.New().String(),
			Username:  config.Envs.INITIAL_USER_USERNAME,
			Password:  hashedPassword,
			Email:     fmt.Sprint(config.Envs.INITIAL_USER_USERNAME, "@gmail.com"),
			Role:      enum.UserRole_User,
			CreatedAt: helper.TimeNowUTC(),
			UpdatedAt: helper.TimeNowUTC(),
		})
	} else {
		logger.Warningf("initial user username and password not set")
	}

	for _, user := range users {
		logger.Infof("seeding user: %s", user.Username)

		// check if user exists
		existing, _ := userRepo.GetByUsername(ctx, user.Username)
		if existing != nil {
			logger.Warningf("user already exists: %s", user.Username)
			continue
		}

		// create user
		err := userRepo.Create(ctx, &user)
		if err != nil {
			logger.Warningf("failed to seed user: %s", user.Username)
			continue
		}

		logger.Infof("user seeded: %s", user.Username)
	}

	return nil
}
