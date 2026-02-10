package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/dto"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/middleware"
	"github.com/GlaciemArgentum/trossage-backend/internal/http/response"
)

// sendMessage sends new message in chat
//
//	@Summary		Send message
//	@Description	Send a new message
//	@Tags			messages
//	@Security		BearerAuth
//	@Param			chat_id	path		int						true	"Chat ID"
//	@Param			request	body		dto.SendMessageRequest	true	"Message body"
//	@Success		201		{object}	dto.SuccessMessageResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid body"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"User is not member"
//	@Failure		404		{object}	dto.ErrorResponse	"Chat not found"
//	@Failure		422		{object}	dto.ErrorResponse	"Message is empty"
//	@Router			/chats/{chat_id}/messages [post]
func (h *handler) sendMessage(ctx *gin.Context) {
	chatID, err := strconv.ParseInt(ctx.Param("chat_id"), 10, 64)
	if err != nil {
		response.HandleError(ctx, derrors.ErrInvalidPathParams)
		return
	}

	var req dto.SendMessageRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	message, err := h.service.CreateMessage(ctx, userID, chatID, req.Text)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.MessageResponse
	resp.Fill(message)

	response.Created(ctx, resp)
}
