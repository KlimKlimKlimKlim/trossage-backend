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

func (r *Repository) InsertUser(ctx context.Context, user models.User) (models.User, error) {
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
			return models.User{}, derrors.ErrUserAlreadyExists
		}

		return models.User{}, err
	}

	return user, nil
}

func (r *Repository) SelectUserByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User

	err := r.db.QueryRow(ctx, users.SelectUserByLoginQuery,
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
		if errors.As(err, &pgx.ErrNoRows) {
			return models.User{}, derrors.ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (r *Repository) SelectUserByID(ctx context.Context, userID int64) (models.User, error) {
	var user models.User

	err := r.db.QueryRow(ctx, users.SelectUserByIDQuery, userID).Scan(
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
			return models.User{}, derrors.ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	err := r.db.QueryRow(ctx, users.UpdateUserQuery,
		user.ID,
		user.DisplayName,
		user.PasswordHash,
	).Scan(
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, derrors.ErrUserNotFound
		}
		return models.User{}, err
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
