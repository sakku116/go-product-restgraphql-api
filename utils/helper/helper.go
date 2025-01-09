package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func ArrayContains(arr interface{}, item interface{}) bool {
	newArr, ok := arr.([]interface{})
	if !ok {
		return false
	}

	for _, v := range newArr {
		if v == item {
			return true
		}
	}
	return false
}

func PrettyJson(data interface{}) string {
	res, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("<failed to parse json: %v>", err.Error())
	}
	return string(res)
}

func TimeNowUTC() time.Time {
	return time.Now().UTC()
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginCtx, ok := ctx.Value("GinContext").(*gin.Context)
	if !ok {
		return nil, errors.New("unable to get Gin context from context")
	}
	return ginCtx, nil
}
