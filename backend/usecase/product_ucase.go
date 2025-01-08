package ucase

import (
	"backend/domain/dto"
	"backend/domain/model"
	"backend/repository"
	error_utils "backend/utils/error"
	"backend/utils/helper"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProductUcase struct {
	productRepo repository.IProductRepo
}

type IProductUcase interface {
	GetByUUID(ctx context.Context, productUUID string) (*dto.GetProductByUUIDRespData, error)
	CreateProduct(
		ctx context.Context,
		currentUser dto.CurrentUser,
		payload dto.CreateProductReq,
	) (*dto.CreateProductRespData, error)
	UpdateProduct(
		ctx context.Context,
		currentUser dto.CurrentUser,
		productUUID string,
		payload dto.UpdateProductReq,
	) (*dto.UpdateProductRespData, error)
	DeleteProduct(
		ctx context.Context,
		currentUser dto.CurrentUser,
		productUUID string,
	) (*dto.DeleteProductRespData, error)
	GetListProduct(
		ctx context.Context,
		params dto.GetProductListReq,
	) (*dto.GetProductListRespData, error)
}

func NewProductUcase(productRepo repository.IProductRepo) IProductUcase {
	return &ProductUcase{productRepo: productRepo}
}

func (ucase *ProductUcase) GetByUUID(ctx context.Context, productUUID string) (*dto.GetProductByUUIDRespData, error) {
	product, err := ucase.productRepo.GetByUUID(ctx, productUUID)
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

	return &dto.GetProductByUUIDRespData{
		ProductModel: *product,
	}, nil
}

func (ucase *ProductUcase) CreateProduct(
	ctx context.Context,
	currentUser dto.CurrentUser,
	payload dto.CreateProductReq,
) (*dto.CreateProductRespData, error) {
	// check if product already exists by current user
	product, _ := ucase.productRepo.GetByName(ctx, payload.Name, &currentUser.UUID)
	logger.Debugf("product by name: %v", product)
	if product != nil {
		logger.Errorf("product with Name %s already exists", payload.Name)
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  fmt.Sprintf("product with Name %s already exists", payload.Name),
		}
	}

	// create product
	product = &model.ProductModel{
		UUID:      uuid.New().String(),
		UserUUID:  currentUser.UUID,
		Name:      payload.Name,
		Price:     payload.Price,
		Stock:     payload.Stock,
		UpdatedAt: helper.TimeNowUTC(),
		CreatedAt: helper.TimeNowUTC(),
	}
	err := product.ValidateBefore()
	if err != nil {
		return nil, err
	}

	err = ucase.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &dto.CreateProductRespData{
		ProductModel: *product,
	}, nil
}

func (ucase *ProductUcase) UpdateProduct(
	ctx context.Context,
	currentUser dto.CurrentUser,
	productUUID string,
	payload dto.UpdateProductReq,
) (*dto.UpdateProductRespData, error) {
	// check existing
	product, err := ucase.productRepo.GetByUUID(ctx, productUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "product not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	// update product obj
	if payload.Name != nil {
		// check if name already exists
		tmp, _ := ucase.productRepo.GetByName(ctx, *payload.Name, &currentUser.UUID)
		if tmp != nil {
			logger.Errorf("product with Name %s already exists", *payload.Name)
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  fmt.Sprintf("product with Name %s already exists", *payload.Name),
			}
		}
		product.Name = *payload.Name
	}
	if payload.Price != nil {
		product.Price = *payload.Price
	}
	if payload.Stock != nil {
		product.Stock = *payload.Stock
	}

	// update product
	err = ucase.productRepo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateProductRespData{
		ProductModel: *product,
	}, nil
}

func (ucase *ProductUcase) DeleteProduct(
	ctx context.Context,
	currentUser dto.CurrentUser,
	productUUID string,
) (*dto.DeleteProductRespData, error) {
	// find product
	product, err := ucase.productRepo.GetByUUID(ctx, productUUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "product not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	// check ownership
	if product.UserUUID != currentUser.UUID {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  "you don't have permission to delete this product",
		}
	}

	// delete product
	err = ucase.productRepo.Delete(ctx, product.UUID)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "product not found",
				Detail:   err.Error(),
			}
		}
		return nil, err
	}

	return &dto.DeleteProductRespData{
		ProductModel: *product,
	}, nil
}

func (ucase *ProductUcase) GetListProduct(
	ctx context.Context,
	params dto.GetProductListReq,
) (*dto.GetProductListRespData, error) {
	products, total, err := ucase.productRepo.GetList(ctx, dto.ProductRepo_GetListParams{
		UserUUID:  params.UserUUID,
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

	resp := &dto.GetProductListRespData{}
	resp.Set(total, *params.Page, *params.Limit)
	resp.Data = products

	return resp, nil
}
