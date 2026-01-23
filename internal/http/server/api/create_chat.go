package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middleware"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// createChat creates a new chat with another user
//
//	@Summary		Create chat
//	@Description	Create a new chat
//	@Tags			chats
//	@Security		BearerAuth
//	@Param			request	body		dto.CreateChatRequest	true	"Chat creation request"
//	@Success		200		{object}	dto.SuccessChatResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		404		{object}	dto.ErrorResponse	"User not found"
//	@Failure		409		{object}	dto.ErrorResponse	"Chat already exists"
//	@Failure		422		{object}	dto.ErrorResponse	"Cannot create chat with yourself"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/chats [post].
func (h *handler) createChat(ctx *gin.Context) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	var req dto.CreateChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	chat, otherUser, err := h.service.CreateChat(ctx, userID, req.UserID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.ChatResponse
	resp.Fill(chat, otherUser)

	response.OK(ctx, resp)
}
