package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// sendTyping sends user's typing operations to other chat participants via WebSocket
//
//	@Summary		Send typing operations
//	@Description	Sends draft message editing operations in real-time.
//	@Description	Operations:
//	@Description	- `insert`: add text at position (requires position, text)
//	@Description	- `delete`: remove characters from position (requires position, length)
//	@Description	- `replace`: replace entire text, e.g. autocorrect (requires position, text)
//	@Description	- `clear`: complete input field clear on send/cancel
//	@Tags			typing
//	@Security		BearerAuth
//	@Param			chat_id	path		int						true	"Chat ID"
//	@Param			request	body		dto.TypingUpdateRequest	true	"Typing operations"
//	@Success		204		{object}	nil						"No Content"
//	@Failure		400		{object}	dto.ErrorResponse		"Invalid body"
//	@Failure		401		{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse		"User is not member"
//	@Failure		404		{object}	dto.ErrorResponse		"Chat not found"
//	@Router			/chats/{chat_id}/typing [post]
func (h *handler) sendTyping(ctx *gin.Context) {
	chatID, err := strconv.ParseInt(ctx.Param("chat_id"), 10, 64)
	if err != nil {
		response.HandleError(ctx, derrors.ErrInvalidPathParams)
		return
	}

	var req dto.TypingUpdateRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	if !req.Validate() {
		response.HandleError(ctx, derrors.ErrInvalidBody)
		return
	}

	userID, ok := middlewares.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	if err = h.service.SendTyping(ctx, userID, chatID, req); err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
