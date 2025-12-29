package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middlewares"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/params"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// getChats returns user's chat list
//
//	@Summary		Get chats
//	@Description	Get user's chat list with pagination
//	@Tags			chats
//	@Security		BearerAuth
//	@Param			limit	query		int	false	"Limit"		default(20)
//	@Param			offset	query		int	false	"Offset"	default(0)
//	@Success		200		{object}	dto.SuccessChatsListResponse
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/chats [get].
func (h *handler) getChats(ctx *gin.Context) {
	limit := params.ParseLimit(ctx.Query("limit"))
	offset := params.ParseOffset(ctx.Query("offset"))

	userID, ok := middlewares.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	chats, total, err := h.service.GetUserChats(ctx, userID, limit, offset)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.ChatsListResponse
	resp.Fill(chats, total, limit, offset)

	response.OK(ctx, resp)
}
