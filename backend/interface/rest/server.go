package rest

import (
	"backend/domain/dto"
	interface_pkg "backend/interface"
	"backend/interface/rest/handler"
	"backend/utils/http_response"

	_ "backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger = logging.MustGetLogger("rest")

func SetupServer(ginEngine *gin.Engine, commonDependencies interface_pkg.CommonDependency) {
	router := ginEngine
	responseWriter := http_response.NewHttpResponseWriter()

	// handlers
	authHandler := handler.NewAuthHandler(responseWriter, commonDependencies.AuthUcase)
	_ = authHandler

	// register routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, dto.BaseJSONResp{
			Code:    200,
			Message: "pong",
		})
	})
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/check-token", authHandler.CheckToken)
	router.POST("/auth/refresh-token", authHandler.RefreshToken)

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/swagger/index.html")
	})
}
