package ucase

import (
	"backend/domain/dto"
	"backend/domain/model"
	"backend/repository"
	bcrypt_util "backend/utils/bcrypt"
	error_utils "backend/utils/error"
	"backend/utils/helper"
	validator_util "backend/utils/validator/user"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserUcase struct {
	userRepo repository.IUserRepo
}

type IUserUcase interface {
	GetByUUID(ctx context.Context, userUUID string) (*dto.GetUserByUUIDResp, error)
	CreateUser(
		ctx context.Context,
		payload dto.CreateUserReq,
	) (*dto.CreateUserRespData, error)
	UpdateUser(
		ctx context.Context,
		userUUID string,
		payload dto.UpdateUserReq,
	) (*dto.UpdateUserRespData, error)
	DeleteUser(
		ctx context.Context,
		userUUID string,
	) (*dto.DeleteUserRespData, error)
	GetUserList(
		ctx context.Context,
		params dto.GetUserListReq,
	) (*dto.GetUserListRespData, error)
}

func NewUserUcase(userRepo repository.IUserRepo) IUserUcase {
	return &UserUcase{userRepo: userRepo}
}

func (ucase *UserUcase) GetByUUID(ctx context.Context, userUUID string) (*dto.GetUserByUUIDResp, error) {
	user, err := ucase.userRepo.GetByUUID(ctx, userUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	return &dto.GetUserByUUIDResp{
		UUID:      user.UUID,
		Username:  user.Username,
		Role:      user.Role.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (ucase *UserUcase) CreateUser(
	ctx context.Context,
	payload dto.CreateUserReq,
) (*dto.CreateUserRespData, error) {
	// validate input
	err := validator_util.ValidateUsername(payload.Username)
	if err != nil {
		logger.Errorf("error validating username: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}

	err = validator_util.ValidateEmail(payload.Email)
	if err != nil {
		logger.Errorf("error validating email: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}

	// check if user exists
	user, _ := ucase.userRepo.GetByEmail(ctx, payload.Email)
	logger.Debugf("user by email: %v", user)
	if user != nil {
		logger.Errorf("user with email %s already exists", payload.Email)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  fmt.Sprintf("user with email %s already exists", payload.Email),
		}
	}

	user, _ = ucase.userRepo.GetByUsername(ctx, payload.Username)
	if user != nil {
		logger.Errorf("user with username %s already exists", payload.Username)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  fmt.Sprintf("user with username %s already exists", payload.Username),
		}
	}

	// create user
	user = &model.UserModel{
		UUID:      uuid.New().String(),
		Username:  payload.Username,
		Password:  payload.Password,
		Role:      "user",
		UpdatedAt: helper.TimeNowUTC(),
		CreatedAt: helper.TimeNowUTC(),
	}
	err = user.ValidateBefore()
	if err != nil {
		return nil, err
	}

	// create password
	user.Password, err = bcrypt_util.Hash(payload.Password)
	if err != nil {
		logger.Errorf("error hashing password: %v", err)
		return nil, err
	}

	err = ucase.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.CreateUserRespData{
		UUID:      user.UUID,
		Username:  user.Username,
		Role:      user.Role.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (ucase *UserUcase) UpdateUser(
	ctx context.Context,
	userUUID string,
	payload dto.UpdateUserReq,
) (*dto.UpdateUserRespData, error) {
	// validate input
	if payload.Username != nil {
		err := validator_util.ValidateUsername(*payload.Username)
		if err != nil {
			logger.Errorf("error validating username: %s", err.Error())
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  err.Error(),
			}
		}
	}

	if payload.Email != nil {
		err := validator_util.ValidateEmail(*payload.Email)
		if err != nil {
			logger.Errorf("error validating email: %s", err.Error())
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  err.Error(),
			}
		}
	}

	if payload.Password != nil {
		err := validator_util.ValidatePassword(*payload.Password)
		if err != nil {
			logger.Errorf("error validating password: %s", err.Error())
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  err.Error(),
			}
		}
	}

	// get existing user
	user, err := ucase.userRepo.GetByUUID(ctx, userUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	// update user obj
	if payload.Username != nil {
		// check if username already exists
		tmp, _ := ucase.userRepo.GetByUsername(ctx, *payload.Username)
		if tmp != nil {
			logger.Errorf("user with username %s already exists", *payload.Username)
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  fmt.Sprintf("user with username %s already exists", *payload.Username),
			}
		}
		user.Username = *payload.Username
	}
	if payload.Email != nil {
		// check if email already exists
		tmp, _ := ucase.userRepo.GetByEmail(ctx, *payload.Email)
		if tmp != nil {
			logger.Errorf("user with email %s already exists", *payload.Email)
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  fmt.Sprintf("user with email %s already exists", *payload.Email),
			}
		}
		user.Email = *payload.Email
	}
	if payload.Password != nil {
		password, err := bcrypt_util.Hash(*payload.Password)
		if err != nil {
			logger.Errorf("error hashing password: %v", err)
			return nil, err
		}
		user.Password = password
	}
	if payload.Role != nil {
		user.Role = *payload.Role
	}

	// update user
	err = ucase.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateUserRespData{
		UUID:      user.UUID,
		Username:  user.Username,
		Role:      user.Role.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (ucase *UserUcase) DeleteUser(
	ctx context.Context,
	userUUID string,
) (*dto.DeleteUserRespData, error) {
	// find user
	user, err := ucase.userRepo.GetByUUID(ctx, userUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	// delete user
	err = ucase.userRepo.Delete(ctx, user.UUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "user not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	return &dto.DeleteUserRespData{
		UUID:      user.UUID,
		Username:  user.Username,
		Role:      user.Role.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, nil
}

func (ucase *UserUcase) GetUserList(
	ctx context.Context,
	params dto.GetUserListReq,
) (*dto.GetUserListRespData, error) {
	users, total, err := ucase.userRepo.GetList(ctx, dto.UserRepo_GetListParams{
		Query:     params.Query,
		QueryBy:   params.QueryBy,
		Page:      params.Page,
		Limit:     params.Limit,
		SortBy:    params.SortBy,
		SortOrder: params.SortOrder,
		DoCount:   true,
	})
	if err != nil {
		return nil, err
	}

	resp := &dto.GetUserListRespData{}
	logger.Debugf("total: %d, page: %d, limit: %d", total, *params.Page, *params.Limit)
	resp.Set(total, *params.Page, *params.Limit)
	resp.Data = users

	return resp, nil
}
