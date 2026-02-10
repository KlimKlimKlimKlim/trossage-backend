package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	derrors "github.com/GlaciemArgentum/trossage-backend/internal/errors"
	"github.com/GlaciemArgentum/trossage-backend/internal/model"
	query "github.com/GlaciemArgentum/trossage-backend/internal/repository/postgres/query/user"
)

func (r *Repository) InsertUser(ctx context.Context, user model.AuthUser) (model.AuthUser, error) {
	err := r.querier.QueryRow(ctx, query.InsertUserQuery,
		user.Login,
		user.PasswordHash,
		user.DisplayName,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return model.AuthUser{}, derrors.ErrUserAlreadyExists
		}

		return model.AuthUser{}, err
	}

	return user, nil
}

func (r *Repository) SelectAuthUserByLogin(ctx context.Context, login string) (model.AuthUser, error) {
	var user model.AuthUser

	err := r.querier.QueryRow(ctx, query.SelectAuthUserByLoginQuery,
		login,
	).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.DisplayName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.AuthUser{}, derrors.ErrUserNotFound
		}

		return model.AuthUser{}, err
	}

	return user, nil
}

func (r *Repository) SelectAuthUserByID(ctx context.Context, userID int64) (model.AuthUser, error) {
	var user model.AuthUser

	err := r.querier.QueryRow(ctx, query.SelectAuthUserByIDQuery, userID).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.DisplayName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.AuthUser{}, derrors.ErrUserNotFound
		}

		return model.AuthUser{}, err
	}

	return user, nil
}

func (r *Repository) SelectUserByID(ctx context.Context, userID int64) (model.User, error) {
	var user model.User

	err := r.querier.QueryRow(ctx, query.SelectUserByIDQuery, userID).Scan(
		&user.ID,
		&user.Login,
		&user.DisplayName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, derrors.ErrUserNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user model.AuthUser) (model.AuthUser, error) {
	err := r.querier.QueryRow(ctx, query.UpdateUserQuery,
		user.ID,
		user.DisplayName,
		user.PasswordHash,
	).Scan(
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.AuthUser{}, derrors.ErrUserNotFound
		}

		return model.AuthUser{}, err
	}

	return user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID int64) error {
	result, err := r.querier.Exec(ctx, query.DeleteUserQuery, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return derrors.ErrUserNotFound
	}

	return nil
}

func (r *Repository) SelectUsersByLoginPrefix(
	ctx context.Context,
	userID int64,
	prefix string,
	limit, offset int,
) ([]model.User, error) {
	rows, err := r.querier.Query(ctx, query.SelectUsersByLoginPrefixQuery, userID, prefix, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.User, 0, limit)

	for rows.Next() {
		var user model.User

		err = rows.Scan(
			&user.ID,
			&user.Login,
			&user.DisplayName,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, rows.Err()
}

func (r *Repository) CountUsersByLoginPrefix(ctx context.Context, userID int64, prefix string) (int, error) {
	var count int
	if err := r.querier.QueryRow(ctx, query.CountUsersByLoginPrefixQuery, userID, prefix).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
