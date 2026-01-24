package postgres

import (
	"context"

	"github.com/jackc/pgerrcode"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
)

func (s *RepositoryTestSuite) TestWithTx() {
	s.Run("creates read-write transaction", func() {
		ctx := context.Background()

		txRepo, commiter, err := s.repo.WithTx(ctx, false)

		s.Require().NoError(err)
		s.NotNil(txRepo)
		s.NotNil(commiter)

		err = commiter.Rollback(ctx)
		s.Require().NoError(err)
	})

	s.Run("creates read-only transaction", func() {
		ctx := context.Background()

		txRepo, commiter, err := s.repo.WithTx(ctx, true)

		s.Require().NoError(err)
		s.NotNil(txRepo)
		s.NotNil(commiter)

		err = commiter.Rollback(ctx)
		s.Require().NoError(err)
	})

	s.Run("commit persists data", func() {
		ctx := context.Background()

		txRepo, commiter, err := s.repo.WithTx(ctx, false)
		s.Require().NoError(err)

		user := model.AuthUser{
			User: model.User{
				Login:       "tx_commit_user",
				DisplayName: "TxCommitUser",
			},
			PasswordHash: "hash_tx_commit",
		}
		insertedUser, err := txRepo.InsertUser(ctx, user)
		s.Require().NoError(err)

		err = commiter.Commit(ctx)
		s.Require().NoError(err)

		foundUser, err := s.repo.SelectAuthUserByID(ctx, insertedUser.ID)
		s.Require().NoError(err)
		s.Equal(insertedUser.ID, foundUser.ID)
		s.Equal(user.Login, foundUser.Login)
	})

	s.Run("rollback discards data", func() {
		ctx := context.Background()

		txRepo, commiter, err := s.repo.WithTx(ctx, false)
		s.Require().NoError(err)

		user := model.AuthUser{
			User: model.User{
				Login:       "tx_rollback_user",
				DisplayName: "TxRollbackUser",
			},
			PasswordHash: "hash_tx_rollback",
		}
		insertedUser, err := txRepo.InsertUser(ctx, user)
		s.Require().NoError(err)
		s.NotZero(insertedUser.ID)

		err = commiter.Rollback(ctx)
		s.Require().NoError(err)

		_, err = s.repo.SelectAuthUserByLogin(ctx, user.Login)
		s.Require().ErrorIs(err, derrors.ErrUserNotFound)
	})

	s.Run("read-only transaction cannot write", func() {
		ctx := context.Background()

		txRepo, commiter, err := s.repo.WithTx(ctx, true)
		s.Require().NoError(err)

		defer func() {
			_ = commiter.Rollback(ctx)
		}()

		user := model.AuthUser{
			User: model.User{
				Login:       "tx_readonly_user",
				DisplayName: "ReadOnlyUser",
			},
			PasswordHash: "hash_readonly",
		}
		_, err = txRepo.InsertUser(ctx, user)

		s.assertPgError(err, pgerrcode.ReadOnlySQLTransaction)
	})

	s.Run("nested transaction not supported", func() {
		ctx := context.Background()

		txRepo, commiter, err := s.repo.WithTx(ctx, false)
		s.Require().NoError(err)

		defer func() {
			_ = commiter.Rollback(ctx)
		}()

		_, _, err = txRepo.WithTx(ctx, false)

		s.Require().ErrorIs(err, derrors.ErrTxNotSupported)
	})

	s.Run("multiple operations in single transaction", func() {
		ctx := context.Background()

		txRepo, commiter, err := s.repo.WithTx(ctx, false)
		s.Require().NoError(err)

		user1 := model.AuthUser{
			User: model.User{
				Login:       "tx_multi_user1",
				DisplayName: "MultiUser1",
			},
			PasswordHash: "hash_multi1",
		}
		insertedUser1, err := txRepo.InsertUser(ctx, user1)
		s.Require().NoError(err)

		user2 := model.AuthUser{
			User: model.User{
				Login:       "tx_multi_user2",
				DisplayName: "MultiUser2",
			},
			PasswordHash: "hash_multi2",
		}
		insertedUser2, err := txRepo.InsertUser(ctx, user2)
		s.Require().NoError(err)

		chat, err := txRepo.InsertChat(ctx, model.ChatTypePrivate)
		s.Require().NoError(err)

		err = txRepo.InsertChatParticipants(ctx, chat.ID, insertedUser1.ID, insertedUser2.ID)
		s.Require().NoError(err)

		err = commiter.Commit(ctx)
		s.Require().NoError(err)

		_, err = s.repo.SelectAuthUserByID(ctx, insertedUser1.ID)
		s.Require().NoError(err)

		_, err = s.repo.SelectAuthUserByID(ctx, insertedUser2.ID)
		s.Require().NoError(err)

		members, err := s.repo.SelectChatMembers(ctx, chat.ID)
		s.Require().NoError(err)
		s.Len(members, 2)
	})
}
