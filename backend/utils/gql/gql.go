package gql_utils

import (
	"backend/domain/dto"
	"backend/domain/enum"
	"backend/utils/helper"
	"context"
	"errors"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("gql_utils")

// !!WARNING!! ginContext must be injected in gin handler to get currentUser through the context
// otherwise it will not be able to retrieve currentUser
func GetCurrentUser(ctx context.Context) (*dto.CurrentUser, error) {
	ginCtx, err := helper.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	currentUser, ok := ginCtx.Value("currentUser").(*dto.CurrentUser)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	return currentUser, nil
}

func GetCurrentUserAdminOnly(ctx context.Context) (*dto.CurrentUser, error) {
	currentUser, err := GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if currentUser.Role != enum.UserRole_Admin {
		return nil, errors.New("forbidden")
	}
	return currentUser, nil
}
