package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/http/dto"
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
func (s *state) searchUsers(ctx *gin.Context) {
	query := strings.ToLower(strings.TrimSpace(ctx.Query("q")))
	if query == "" {
		response.HandleError(ctx, derrors.ErrEmptyQuery)
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.controller.SearchUsersByLogin(ctx, query, limit, offset)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	var resp dto.UsersSearchResponse
	resp.Fill(users, total, limit, offset)

	response.OK(ctx, resp)
}
