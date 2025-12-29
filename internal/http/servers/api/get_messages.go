package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/params"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// getMessages gets messages from chat
//
//	@Summary		Get messages
//	@Description	Get messages from chat
//	@Tags			messages
//	@Security		BearerAuth
//	@Param			chat_id	path		int	true	"Chat ID"
//	@Param			limit	query		int	false	"Limit"		default(20)
//	@Param			offset	query		int	false	"Offset"	default(0)
//	@Success		200		{object}	dto.SuccessMessagesResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid path params"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"User is not member"
//	@Failure		404		{object}	dto.ErrorResponse	"Chat not found"
//	@Router			/chats/{chat_id}/messages [get]
func (h *handler) getMessages(ctx *gin.Context) {
	chatID, err := strconv.ParseInt(ctx.Param("chat_id"), 10, 64)
	if err != nil {
		response.HandleError(ctx, derrors.ErrInvalidPathParams)
		return
	}

	limit := params.ParseLimit(ctx.Query("limit"))
	offset := params.ParseOffset(ctx.Query("offset"))

	userID, ok := middlewares.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	messages, total, err := h.service.GetMessages(ctx, userID, chatID, limit, offset)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.MessagesResponse
	resp.Fill(messages, limit, offset, total)

	response.Created(ctx, resp)
}
