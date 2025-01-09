package dto

import (
	"backend/domain/enum"
	"backend/domain/model"
	"errors"
	"time"
)

type GetUserByUUIDResp struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
}

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

type CreateUserRespData struct {
	UUID      string    `json:"uuid"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserReq struct {
	Username *string        `json:"username"`
	Email    *string        `json:"email"`
	Password *string        `json:"password"`
	Role     *enum.UserRole `json:"role" binding:"oneof=admin user"`
}

type UpdateUserRespData struct {
	UUID      string    `json:"uuid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteUserRespData struct {
	UUID      string    `json:"uuid"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserListReq struct {
	Query     *string `form:"query" binding:"omitempty"`
	QueryBy   *string `form:"query_by" binding:"omitempty,oneof=username email role" default:"username"`
	Page      *int    `form:"page" binding:"omitempty,gt=0" default:"1"`
	Limit     *int    `form:"limit" binding:"omitempty,gt=0" default:"10"`
	SortOrder *int    `form:"sort_order" binding:"omitempty,oneof=1 -1" default:"-1"`
	SortBy    *string `form:"sort_by" binding:"omitempty,oneof=updated_at created_at username email role" default:"updated_at"`
}

type GetUserListRespData struct {
	BasePaginationRespData
	Data []model.UserModel `json:"data"`
}

type UserRepo_GetListParams struct {
	Query     *string
	QueryBy   *string
	Page      *int
	Limit     *int
	SortOrder *int    `default:"-1"`
	SortBy    *string `default:"updated_at"`
	DoCount   bool
}

func (dto *UserRepo_GetListParams) Validate() error {
	if dto.QueryBy != nil {
		if *dto.QueryBy != "username" &&
			*dto.QueryBy != "email" &&
			*dto.QueryBy != "role" {
			return errors.New("invalid query by")
		}
	}

	if dto.SortBy != nil {
		if *dto.SortBy != "updated_at" &&
			*dto.SortBy != "created_at" &&
			*dto.SortBy != "username" &&
			*dto.SortBy != "email" &&
			*dto.SortBy != "role" {
			return errors.New("invalid sort by")
		}
	}

	if dto.SortOrder != nil {
		if *dto.SortOrder != 1 && *dto.SortOrder != -1 {
			return errors.New("invalid sort order")
		}
	}

	return nil
}
