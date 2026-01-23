package postgres

import (
	"context"
	"time"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

type IRepository interface { //nolint:interfacebloat // it's okay for repository interfaces
	InsertUser(ctx context.Context, user model.AuthUser) (model.AuthUser, error)
	SelectAuthUserByLogin(ctx context.Context, login string) (model.AuthUser, error)
	SelectAuthUserByID(ctx context.Context, userID int64) (model.AuthUser, error)
	SelectUserByID(ctx context.Context, userID int64) (model.User, error)
	UpdateUser(ctx context.Context, user model.AuthUser) (model.AuthUser, error)
	DeleteUser(ctx context.Context, userID int64) error
	SelectUsersByLoginPrefix(ctx context.Context, userID int64, prefix string, limit, offset int) ([]model.User, error)
	CountUsersByLoginPrefix(ctx context.Context, userID int64, prefix string) (int, error)

	InsertRefreshToken(ctx context.Context, token model.Token) (model.Token, error)
	SelectRefreshToken(ctx context.Context, userID int64, tokenHash string) (model.Token, error)
	RevokeRefreshTokenByID(ctx context.Context, tokenID int64) error
	RevokeRefreshTokensByUserID(ctx context.Context, userID int64) error
	DeleteExpiredTokens(ctx context.Context, olderThan time.Time) (int64, error)
	DeleteRevokedTokens(ctx context.Context, olderThan time.Time) (int64, error)

	SelectChatBetweenUsers(ctx context.Context, userID1, userID2 int64) (int64, error)
	InsertChat(ctx context.Context, chatType model.ChatType) (model.Chat, error)
	InsertChatParticipants(ctx context.Context, chatID int64, userIDs ...int64) error
	SelectUserChats(ctx context.Context, userID int64, limit, offset int) ([]model.ChatWithDetails, error)
	CountUserChats(ctx context.Context, userID int64) (int, error)
	IsUserMember(ctx context.Context, chatID, userID int64) (bool, error)
	SelectChatByID(ctx context.Context, chatID int64) (model.Chat, error)
	SelectChatMembers(ctx context.Context, chatID int64) ([]int64, error)

	CreateMessage(ctx context.Context, chatID, senderID int64, text string) (model.Message, error)
	SelectMessages(ctx context.Context, chatID int64, limit, offset int) ([]model.Message, error)
	CountChatMessages(ctx context.Context, chatID int64) (int, error)
}
