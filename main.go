package main

import (
	"backend/config"
	"backend/domain/model"
	interface_pkg "backend/interface"
	"backend/interface/gql"
	"backend/interface/rest"
	"backend/repository"
	ucase "backend/usecase"
	"backend/utils/helper"
	mongo_util "backend/utils/mongo"
	seeder_util "backend/utils/seeder/user"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	productModel := model.ProductModel{}

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
	productRepo := repository.NewProductRepo(
		mongoDatabase.Collection(
			productModel.GetMongoProps().CollName,
		),
	)

	// usecases
	authUcase := ucase.NewAuthUcase(userRepo, refreshTokenRepo)
	userUcase := ucase.NewUserUcase(userRepo)
	productUcase := ucase.NewProductUcase(productRepo)

	dependencies := interface_pkg.CommonDependency{
		AuthUcase:    authUcase,
		UserUcase:    userUcase,
		ProductUcase: productUcase,
	}

	// proccess args
	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "--seed-users":
			seeder_util.SeedUser(userRepo)
		case "--ensure-indexes":
			mongo_util.EnsureMongoIndexes(
				mongoDatabase,
				userModel,
				refreshTokenModel,
				productModel,
			)
		}
	}

	// rest + gql
	ginEngine := gin.Default()
	gql.SetupGql(ginEngine, dependencies)
	rest.SetupRest(ginEngine, dependencies)

	// http server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Envs.HOST, config.Envs.PORT),
		Handler: ginEngine,
	}

	// run server
	go func() {
		logger.Debugf("starting server at %s:%d", config.Envs.HOST, config.Envs.PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	// gracefull shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("server forced to shutdown: %v", err)
	}

	logger.Info("server exiting")
}
