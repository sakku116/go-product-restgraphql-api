package dto

import "backend/domain/enum"

type CurrentUser struct {
	UUID     string        `json:"uuid"`
	Username string        `json:"username"`
	Role     enum.UserRole `json:"role"`
	Email    string        `json:"email"`
}

type RegisterUserReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type RegisterUserRespData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRespData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CheckTokenReq struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type CheckTokenRespData struct {
	UUID     string        `json:"uuid"`
	Username string        `json:"username"`
	Role     enum.UserRole `json:"role"`
	Email    string        `json:"email"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenRespData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
