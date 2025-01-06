package rest

import (
	"backend/domain/dto"
	interface_pkg "backend/interface"
	"backend/interface/rest/handler"
	"backend/middleware"
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
	userHandler := handler.NewUserHandler(responseWriter, commonDependencies.UserUcase)

	// middlewares
	authMiddleware := middleware.AuthMiddleware(responseWriter)
	adminOnlyMiddleware := middleware.AuthAdminOnlyMiddleware(responseWriter)

	// register routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, dto.BaseJSONResp{
			Code:    200,
			Message: "pong",
		})
	})

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", authHandler.Register)
		authRouter.POST("/login", authHandler.Login)
		authRouter.POST("/check-token", authHandler.CheckToken)
		authRouter.POST("/refresh-token", authHandler.RefreshToken)
	}

	secureRouter := router.Group("/", authMiddleware)
	{
		userRouter := secureRouter.Group("/users")
		{
			userRouter.GET("/me", userHandler.GetUserMe)
			userRouter.GET("/:uuid", userHandler.GetUserByUUID)
			userRouter.PUT("/me", userHandler.UpdateUserMe)

			// admin only
			userRouter.POST("", adminOnlyMiddleware, userHandler.CreateUser)
			userRouter.PUT("/:uuid", adminOnlyMiddleware, userHandler.UpdateUser)
			userRouter.DELETE("/:uuid", adminOnlyMiddleware, userHandler.DeleteUser)
		}
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/swagger/index.html")
	})
}
