package api

import (
	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/middleware"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/params"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/response"
)

// searchUsers returns users matching login prefix
//
//	@Summary		Search users
//	@Description	Search users by login prefix with pagination
//	@Tags			users
//	@Security		BearerAuth
//	@Param			q		query		string	true	"Search query (login prefix)"
//	@Param			limit	query		int		false	"Limit"		default(20)
//	@Param			offset	query		int		false	"Offset"	default(0)
//	@Success		200		{object}	dto.SuccessUsersSearchResponse
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		422		{object}	dto.ErrorResponse	"Invalid query"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/users/search [get]
func (h *handler) searchUsers(ctx *gin.Context) {
	loginPrefix, err := params.ParseSearchQuery(ctx.Query("q"))
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	limit := params.ParseLimit(ctx.Query("limit"))
	offset := params.ParseOffset(ctx.Query("offset"))

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		response.HandleError(ctx, derrors.ErrUserIDIsEmpty)
		return
	}

	users, total, err := h.service.SearchUsersByLogin(ctx, userID, loginPrefix, limit, offset)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.UsersSearchResponse
	resp.Fill(users, total, limit, offset)

	response.OK(ctx, resp)
}
