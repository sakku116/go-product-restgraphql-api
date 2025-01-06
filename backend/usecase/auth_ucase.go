package ucase

import (
	"backend/config"
	"backend/domain/dto"
	"backend/domain/enum"
	"backend/domain/model"
	"backend/repository"
	bcrypt_util "backend/utils/bcrypt"
	error_utils "backend/utils/error"
	"backend/utils/helper"
	jwt_util "backend/utils/jwt"
	validator_util "backend/utils/validator/user"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthUcase struct {
	userRepo         repository.IUserRepo
	refreshTokenRepo repository.IRefreshTokenRepo
}

type IAuthUcase interface {
	Register(ctx context.Context, payload dto.RegisterUserReq) (*dto.RegisterUserRespData, error)
	Login(ctx context.Context, payload dto.LoginReq) (*dto.LoginRespData, error)
	RefreshToken(ctx context.Context, payload dto.RefreshTokenReq) (*dto.RefreshTokenRespData, error)
	CheckToken(payload dto.CheckTokenReq) (*dto.CheckTokenRespData, error)
}

func NewAuthUcase(
	userRepo repository.IUserRepo,
	refreshTokenRepo repository.IRefreshTokenRepo,
) IAuthUcase {
	return &AuthUcase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *AuthUcase) Register(ctx context.Context, payload dto.RegisterUserReq) (*dto.RegisterUserRespData, error) {
	// check if user exists
	err := validator_util.ValidateUsername(payload.Username)
	if err != nil {
		logger.Errorf("error validating username: %s", err.Error())
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}
	user, _ := s.userRepo.GetByUsername(ctx, payload.Username)
	logger.Debugf("user by username: %v", user)
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
		Email:     payload.Email,
		Role:      enum.UserRole_User,
		CreatedAt: helper.TimeNowUTC(),
		UpdatedAt: helper.TimeNowUTC(),
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

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// generate token
	token, err := jwt_util.GenerateJwtToken(user, config.Envs.JWT_SECRET_KEY, config.Envs.JWT_EXP_MINS, nil)
	if err != nil {
		logger.Errorf("error generating token: %v", err)
		return nil, err
	}

	// invalidate old refresh token
	s.refreshTokenRepo.InvalidateManyByUserUUID(ctx, user.UUID)

	// create refresh token
	refreshTokenExpiredAt := helper.TimeNowUTC().Add(time.Minute * time.Duration(config.Envs.REFRESH_TOKEN_EXP_MINS))
	newRefreshTokenObj := model.RefreshTokenModel{
		Token:     uuid.New().String(),
		UserUUID:  user.UUID,
		UsedAt:    nil,
		ExpiredAt: &refreshTokenExpiredAt,
		CreatedAt: helper.TimeNowUTC(),
		UpdatedAt: helper.TimeNowUTC(),
	}
	logger.Debugf("new refresh token: %+v", newRefreshTokenObj)
	err = s.refreshTokenRepo.Create(ctx, &newRefreshTokenObj)
	if err != nil {
		logger.Errorf("error creating refresh token: %v", err)
		return nil, err
	}

	resp := &dto.RegisterUserRespData{
		AccessToken:  token,
		RefreshToken: newRefreshTokenObj.Token,
	}
	return resp, nil
}

func (s *AuthUcase) Login(ctx context.Context, payload dto.LoginReq) (*dto.LoginRespData, error) {
	logger.Debugf("payload: %v", payload)
	// validate username
	err := validator_util.ValidateUsername(payload.Username)
	if err != nil {
		logger.Errorf("invalid username: %s\n%v", payload.Username, err)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}

	// validate password
	err = validator_util.ValidatePassword(payload.Password)
	if err != nil {
		logger.Errorf("invalid password: %s\n%v", payload.Password, err)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  err.Error(),
		}
	}

	// check if user exists
	var existing_user *model.UserModel
	existing_user, _ = s.userRepo.GetByUsername(ctx, payload.Username)
	if existing_user == nil {
		logger.Errorf("user not found")
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Credentials",
		}
	}
	logger.Debugf("user by username: %v", helper.PrettyJson(existing_user))

	// check password
	if !bcrypt_util.Compare(payload.Password, existing_user.Password) {
		logger.Errorf("invalid password")
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Credentials",
		}
	}

	// generate token
	token, err := jwt_util.GenerateJwtToken(existing_user, config.Envs.JWT_SECRET_KEY, config.Envs.JWT_EXP_MINS, nil)
	if err != nil {
		logger.Errorf("error generating token: %v", err)
		return nil, err
	}

	// invalidate old refresh token
	err = s.refreshTokenRepo.InvalidateManyByUserUUID(ctx, existing_user.UUID)
	if err != nil {
		logger.Errorf("error invalidating old refresh token: %v", err)
		return nil, err
	}

	// create refresh token
	refreshTokenExpiredAt := helper.TimeNowUTC().Add(time.Minute * time.Duration(config.Envs.REFRESH_TOKEN_EXP_MINS))
	newRefreshTokenObj := model.RefreshTokenModel{
		Token:     uuid.New().String(),
		UserUUID:  existing_user.UUID,
		UsedAt:    nil,
		ExpiredAt: &refreshTokenExpiredAt,
		CreatedAt: helper.TimeNowUTC(),
		UpdatedAt: helper.TimeNowUTC(),
	}
	logger.Debugf("new refresh token: %+v", helper.PrettyJson(newRefreshTokenObj))
	err = s.refreshTokenRepo.Create(ctx, &newRefreshTokenObj)
	if err != nil {
		logger.Errorf("error creating refresh token: %v", err)
		return nil, err
	}

	return &dto.LoginRespData{
		AccessToken:  token,
		RefreshToken: newRefreshTokenObj.Token,
	}, nil
}

func (s *AuthUcase) RefreshToken(ctx context.Context, payload dto.RefreshTokenReq) (*dto.RefreshTokenRespData, error) {
	// get refresh token
	refreshToken, err := s.refreshTokenRepo.GetByToken(ctx, payload.RefreshToken)
	if err != nil {
		logger.Errorf("refresh token not found: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Refresh Token",
		}
	}

	// check if refresh token is expired
	if refreshToken.ExpiredAt != nil {
		if refreshToken.ExpiredAt.Before(helper.TimeNowUTC()) {
			logger.Errorf("refresh token is expired")
			return nil, &error_utils.CustomErr{
				HttpCode: 401,
				Message:  "Invalid Refresh Token",
			}
		}
	}

	// check if refresh token is used
	if refreshToken.UsedAt != nil {
		logger.Errorf("refresh token is used")
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Refresh Token",
		}
	}

	// check if refresh token is valid
	if refreshToken.Invalid {
		logger.Errorf("refresh token is invalid")
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Refresh Token",
		}
	}

	// mark refresh token as used
	timeNow := helper.TimeNowUTC()
	refreshToken.UsedAt = &timeNow
	refreshToken.UpdatedAt = timeNow
	err = s.refreshTokenRepo.Update(ctx, refreshToken)
	if err != nil {
		logger.Errorf("error updating refresh token: %v", err)
		return nil, err
	}

	// get user
	logger.Debug("refresh token: %+v", helper.PrettyJson(refreshToken))
	user, err := s.userRepo.GetByUUID(ctx, refreshToken.UserUUID)
	if err != nil {
		logger.Errorf("user not found: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "Internal server error",
			Detail:   err.Error(),
		}
	}

	// generate token
	token, err := jwt_util.GenerateJwtToken(user, config.Envs.JWT_SECRET_KEY, config.Envs.JWT_EXP_MINS, nil)
	if err != nil {
		logger.Errorf("error generating token: %v", err)
		return nil, err
	}

	// invalidate old refresh token
	err = s.refreshTokenRepo.InvalidateManyByUserUUID(ctx, user.UUID)
	if err != nil {
		logger.Errorf("error invalidating old refresh token: %v", err)
		return nil, err
	}

	// create refresh token
	refreshTokenExpiredAt := helper.TimeNowUTC().Add(time.Minute * time.Duration(config.Envs.REFRESH_TOKEN_EXP_MINS))
	newRefreshTokenObj := model.RefreshTokenModel{
		Token:     uuid.New().String(),
		UserUUID:  user.UUID,
		UsedAt:    nil,
		ExpiredAt: &refreshTokenExpiredAt,
		CreatedAt: helper.TimeNowUTC(),
		UpdatedAt: helper.TimeNowUTC(),
	}
	err = s.refreshTokenRepo.Create(ctx, &newRefreshTokenObj)
	if err != nil {
		logger.Errorf("error creating refresh token: %v", err)
		return nil, err
	}

	return &dto.RefreshTokenRespData{
		AccessToken:  token,
		RefreshToken: newRefreshTokenObj.Token,
	}, nil
}

func (s *AuthUcase) CheckToken(payload dto.CheckTokenReq) (*dto.CheckTokenRespData, error) {
	claims, err := jwt_util.ValidateJWT(payload.AccessToken, config.Envs.JWT_SECRET_KEY)
	if err != nil || claims == nil {
		logger.Errorf("error validating token: %v", err)
		return nil, &error_utils.CustomErr{
			HttpCode: 401,
			Message:  "Invalid Access Token",
			Detail:   err.Error(),
		}
	}

	resp := &dto.CheckTokenRespData{
		UUID:     claims.UUID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}

	return resp, nil
}
