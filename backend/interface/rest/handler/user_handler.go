package handler

import (
	"backend/domain/dto"
	ucase "backend/usecase"
	"backend/utils/http_response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	respWriter http_response.IHttpResponseWriter
	userUcase  ucase.IUserUcase
}

type IUserHandler interface {
	GetUserByUUID(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

func NewUserHandler(respWriter http_response.IHttpResponseWriter, userUcase ucase.IUserUcase) UserHandler {
	return UserHandler{
		respWriter: respWriter,
		userUcase:  userUcase,
	}
}

// @Summary get user (me)
// @Tags User
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetUserByUUIDResp}
// @Router /users/me [get]
// @Security BearerAuth
func (h *UserHandler) GetUserMe(ctx *gin.Context) {
	currentUser, ok := ctx.MustGet("currentUser").(*dto.CurrentUser)
	if !ok {
		h.respWriter.HTTPJson(ctx, 500, "internal server error", "current user not found", nil)
		return
	}
	userUUID := currentUser.UUID

	data, err := h.userUcase.GetByUUID(ctx, userUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary get user by uuid
// @Tags User
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetUserByUUIDResp}
// @Router /users/{uuid} [get]
// @param uuid path string true "user uuid"
// @Security BearerAuth
func (h *UserHandler) GetUserByUUID(ctx *gin.Context) {
	userUUID := ctx.Param("uuid")

	data, err := h.userUcase.GetByUUID(ctx, userUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary create new user (admin only)
// @Tags User
// @Success 200 {object} dto.BaseJSONResp{data=dto.CreateUserRespData}
// @Router /users [post]
// @param payload  body  dto.CreateUserReq  true "payload"
// @Security BearerAuth
func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var payload dto.CreateUserReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		// h.respWriter.HTTcPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	data, err := h.userUcase.CreateUser(ctx, payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary update user (me)
// @Tags User
// @Success 200 {object} dto.BaseJSONResp{data=dto.UpdateUserRespData}
// @Router /users/me [put]
// @param payload  body  dto.UpdateUserReq  true "payload"
// @Security BearerAuth
func (h *UserHandler) UpdateUserMe(ctx *gin.Context) {
	// get current user
	currentUser, ok := ctx.MustGet("currentUser").(dto.CurrentUser)
	if !ok {
		h.respWriter.HTTPJson(ctx, 500, "internal server error", "current user not found", nil)
		return
	}
	userUUID := currentUser.UUID

	var payload dto.UpdateUserReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	data, err := h.userUcase.UpdateUser(ctx, userUUID, payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary update user (admin only)
// @Tags User
// @Success 200 {object} dto.BaseJSONResp{data=dto.UpdateUserRespData}
// @Router /users/{uuid} [put]
// @param uuid path string true "user uuid"
// @param payload  body  dto.UpdateUserReq  true "payload"
// @Security BearerAuth
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	userUUID := ctx.Param("uuid")
	var payload dto.UpdateUserReq
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		h.respWriter.HTTPJson(ctx, 400, "invalid payload", err.Error(), nil)
		return
	}

	data, err := h.userUcase.UpdateUser(ctx, userUUID, payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}

// @Summary delete user (admin only)
// @Tags User
// @Success 200 {object} dto.BaseJSONResp{data=dto.DeleteUserRespData}
// @Router /users/{uuid} [delete]
// @param uuid path string true "user uuid"
// @Security BearerAuth
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	userUUID := ctx.Param("uuid")

	data, err := h.userUcase.DeleteUser(ctx, userUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(ctx, err)
		return
	}

	h.respWriter.HTTPJsonOK(ctx, data)
}
