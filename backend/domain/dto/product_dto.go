package dto

import (
	"backend/domain/model"
	"errors"
)

type GetProductByUUIDRespData struct {
	model.ProductModel
}

type CreateProductReq struct {
	Name  string  `json:"name" binding:"required"`
	Stock int64   `json:"stock" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type CreateProductRespData struct {
	model.ProductModel
}

type UpdateProductReq struct {
	Name  *string  `json:"name"`
	Stock *int64   `json:"stock"`
	Price *float64 `json:"price"`
}

type UpdateProductRespData struct {
	model.ProductModel
}

type DeleteProductRespData struct {
	model.ProductModel
}

type ProductRepo_GetListParams struct {
	UserUUID  *string
	Query     *string
	QueryBy   *string
	Page      *int
	Limit     *int
	SortOrder *int    `default:"-1"`
	SortBy    *string `default:"updated_at"`
	DoCount   bool
}

func (dto *ProductRepo_GetListParams) Validate() error {
	if dto.QueryBy != nil {
		if *dto.QueryBy != "name" {
			return errors.New("invalid query by")
		}
	}

	if dto.SortBy != nil {
		if *dto.SortBy != "updated_at" &&
			*dto.SortBy != "created_at" &&
			*dto.SortBy != "name" &&
			*dto.SortBy != "price" {
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

type GetProductListReq struct {
	UserUUID  *string `form:"user_uuid" binding:"omitempty"`
	Query     *string `form:"query" binding:"omitempty"`
	QueryBy   *string `form:"query_by" binding:"omitempty,oneof=name" default:"name"`
	Page      *int    `form:"page" binding:"omitempty,gt=0" default:"1"`
	Limit     *int    `form:"limit" binding:"omitempty,gt=0" default:"10"`
	SortOrder *int    `form:"sort_order" binding:"omitempty,oneof=1 -1" default:"-1"`
	SortBy    *string `form:"sort_by" binding:"omitempty,oneof=updated_at created_at name price" default:"updated_at"`
}

type GetProductListRespData struct {
	BasePaginationRespData
	Data []model.ProductModel `json:"data"`
}
