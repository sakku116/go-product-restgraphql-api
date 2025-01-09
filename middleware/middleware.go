package middleware

import (
	"backend/config"
	"backend/domain/dto"
	"backend/utils/http_response"
	jwt_util "backend/utils/jwt"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

func AuthMiddleware(respWriter http_response.IHttpResponseWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debugf("auth middleware")
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			respWriter.HTTPJson(
				c, 401, "unauthorized", "invalid token", nil,
			)
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		currentUser, err := jwt_util.ValidateJWT(token, config.Envs.JWT_SECRET_KEY)
		if err != nil {
			respWriter.HTTPJson(
				c, 401, "unauthorized", err.Error(), nil,
			)
			c.Abort()
			return
		}

		c.Set("currentUser", currentUser)
		c.Next()
	}
}

// need to combine with AuthMiddleware first
func AuthAdminOnlyMiddleware(respWriter http_response.IHttpResponseWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debugf("auth admin only middleware")
		currentUserRaw, ok := c.Get("currentUser")
		if !ok {
			respWriter.HTTPJson(
				c, 500, "internal service error", "current user not found", nil,
			)
			c.Abort()
			return
		}

		currentUser, ok := currentUserRaw.(*dto.CurrentUser)
		if !ok {
			respWriter.HTTPJson(
				c, 500, "internal service error", "current user missmatched", nil,
			)
			c.Abort()
			return
		}

		if currentUser.Role != "admin" {
			respWriter.HTTPJson(
				c, 403, "forbidden", "admin only", nil,
			)
			c.Abort()
			return
		}

		c.Next()
	}
}

func GinContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "GinContext", c))
		c.Next()
	}
}
