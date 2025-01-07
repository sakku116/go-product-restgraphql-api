package handler

import (
	"backend/domain/dto"
	ucase "backend/usecase"
	"backend/utils/http_response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	respWriter   http_response.IHttpResponseWriter
	productUcase ucase.IProductUcase
}

type IProductHandler interface {
	GetProductByUUID(ctx *gin.Context)
	CreateProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	GetProductList(ctx *gin.Context)
}

func NewProductHandler(respWriter http_response.IHttpResponseWriter, productUcase ucase.IProductUcase) IProductHandler {
	return &ProductHandler{
		respWriter:   respWriter,
		productUcase: productUcase,
	}
}

// @Summary get product by uuid
// @Tags Product
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetProductByUUIDRespData}
// @Router /products/{uuid} [get]
// @param uuid path string true "user uuid"
// @Security BearerAuth
func (h *ProductHandler) GetProductByUUID(ctx *gin.Context) {
	userUUID := ctx.Param("uuid")

	data, err := h.productUcase.GetByUUID(ctx, userUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary create new product
// @Tags Product
// @Success 200 {object} dto.BaseJSONResp{data=dto.CreateProductRespData}
// @Router /products [post]
// @param payload  body  dto.CreateProductReq  true "payload"
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var payload dto.CreateProductReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		// h.respWriter.HTTcPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	// get current user
	currentUser, ok := ctx.MustGet("currentUser").(*dto.CurrentUser)
	if !ok {
		h.respWriter.HTTPJson(ctx, 500, "internal server error", "current user not found", nil)
		return
	}

	data, err := h.productUcase.CreateProduct(ctx, *currentUser, payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary update product
// @Tags Product
// @Success 200 {object} dto.BaseJSONResp{data=dto.UpdateProductRespData}
// @Router /products/{uuid} [put]
// @param uuid path string true "product uuid"
// @param payload  body  dto.UpdateProductReq  true "payload"
// @Security BearerAuth
func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	productUUID := ctx.Param("uuid")
	var payload dto.UpdateProductReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	// get current user
	currentUser, ok := ctx.MustGet("currentUser").(*dto.CurrentUser)
	if !ok {
		h.respWriter.HTTPJson(ctx, 500, "internal server error", "current user not found", nil)
		return
	}

	data, err := h.productUcase.UpdateProduct(ctx, *currentUser, productUUID, payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary delete product
// @Tags Product
// @Success 200 {object} dto.BaseJSONResp{data=dto.DeleteProductRespData}
// @Router /products/{uuid} [delete]
// @param uuid path string true "product uuid"
// @Security BearerAuth
func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productUUID := ctx.Param("uuid")

	// get current user
	currentUser, ok := ctx.MustGet("currentUser").(*dto.CurrentUser)
	if !ok {
		h.respWriter.HTTPJson(ctx, 500, "internal server error", "current user not found", nil)
		return
	}

	data, err := h.productUcase.DeleteProduct(ctx, *currentUser, productUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary get product list
// @Tags Product
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetProductListRespData}
// @Router /products [get]
// @param payload  body  dto.GetProductListReq  true "payload"
// @Security BearerAuth
func (h ProductHandler) GetProductList(ctx *gin.Context) {
	var params dto.GetProductListReq
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		h.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	data, err := h.productUcase.GetListProduct(ctx, params)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}
