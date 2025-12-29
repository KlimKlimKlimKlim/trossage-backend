package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/models"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/postgres/repository/queries/users"
)

func (r *Repository) InsertUser(ctx context.Context, user models.AuthUser) (models.AuthUser, error) {
	err := r.db.QueryRow(ctx, users.InsertUserQuery,
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
			return models.AuthUser{}, derrors.ErrUserAlreadyExists
		}

		return models.AuthUser{}, err
	}

	return user, nil
}

func (r *Repository) SelectAuthUserByLogin(ctx context.Context, login string) (models.AuthUser, error) {
	var user models.AuthUser

	err := r.db.QueryRow(ctx, users.SelectAuthUserByLoginQuery,
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
			return models.AuthUser{}, derrors.ErrUserNotFound
		}

		return models.AuthUser{}, err
	}

	return user, nil
}

func (r *Repository) SelectAuthUserByID(ctx context.Context, userID int64) (models.AuthUser, error) {
	var user models.AuthUser

	err := r.db.QueryRow(ctx, users.SelectAuthUserByIDQuery, userID).Scan(
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
			return models.AuthUser{}, derrors.ErrUserNotFound
		}

		return models.AuthUser{}, err
	}

	return user, nil
}

func (r *Repository) SelectUserByID(ctx context.Context, userID int64) (models.User, error) {
	var user models.User

	err := r.db.QueryRow(ctx, users.SelectUserByIDQuery, userID).Scan(
		&user.ID,
		&user.Login,
		&user.DisplayName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, derrors.ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user models.AuthUser) (models.AuthUser, error) {
	err := r.db.QueryRow(ctx, users.UpdateUserQuery,
		user.ID,
		user.DisplayName,
		user.PasswordHash,
	).Scan(
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.AuthUser{}, derrors.ErrUserNotFound
		}

		return models.AuthUser{}, err
	}

	return user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID int64) error {
	result, err := r.db.Exec(ctx, users.DeleteUserQuery, userID)
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
	query string,
	limit, offset int,
) ([]models.User, error) {
	rows, err := r.db.Query(ctx, users.SelectUsersByLoginPrefixQuery, userID, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.User, 0, limit)

	for rows.Next() {
		var user models.User

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

func (r *Repository) CountUsersByLoginPrefix(ctx context.Context, userID int64, query string) (int, error) {
	var count int
	if err := r.db.QueryRow(ctx, users.CountUsersByLoginPrefixQuery, userID, query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
