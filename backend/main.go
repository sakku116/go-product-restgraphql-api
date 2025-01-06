package main

import (
	"backend/config"
	"backend/domain/model"
	interface_pkg "backend/interface"
	"backend/interface/rest"
	"backend/repository"
	ucase "backend/usecase"
	"backend/utils/helper"
	seeder_util "backend/utils/seeder/user"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

func init() {
	config.InitEnv("./.env")
	config.ConfigureLogger()
}

var logger = logging.MustGetLogger("main")

// @title Auth Service RESTful API
// @securitydefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme (add 'Bearer ' prefix).
func main() {
	logger.Debugf("Envs: %v", helper.PrettyJson(config.Envs))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoConn := config.NewMongoConn(ctx)
	mongoDatabase := mongoConn.Database(config.Envs.MONGO_DB)

	// models
	userModel := model.UserModel{}
	refreshTokenModel := model.RefreshTokenModel{}

	// repositories
	userRepo := repository.NewUserRepo(
		mongoDatabase.Collection(
			userModel.GetMongoProps().CollName,
		),
	)
	refreshTokenRepo := repository.NewRefreshTokenRepo(
		mongoDatabase.Collection(
			refreshTokenModel.GetMongoProps().CollName,
		),
	)

	// usecases
	authUcase := ucase.NewAuthUcase(userRepo, refreshTokenRepo)
	userUcase := ucase.NewUserUcase(userRepo)

	dependencies := interface_pkg.CommonDependency{
		AuthUcase: authUcase,
		UserUcase: userUcase,
	}

	// seed data
	seeder_util.SeedUser(userRepo)

	// rest
	ginEngine := gin.Default()
	rest.SetupServer(ginEngine, dependencies)
	go func() {
		logger.Debugf("starting server at %s:%d", config.Envs.HOST, config.Envs.PORT)
		ginEngine.Run(fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT))
	}()

	select {}
}
