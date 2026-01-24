package postgres

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/KlimKlimKlimKlim/trossage-backend/internal/config"
	"github.com/KlimKlimKlimKlim/trossage-backend/internal/model"
	"github.com/KlimKlimKlimKlim/trossage-backend/migrations"
)

const (
	containerLogOccurrences = 2
	containerStartupTimeout = 60 * time.Second
)

type RepositoryTestSuite struct {
	suite.Suite

	container *postgres.PostgresContainer
	pool      *pgxpool.Pool
	repo      *Repository
}

func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	cfg, err := config.New("../../../.env.example")
	s.Require().NoError(err, "failed to load config")

	s.container, err = postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(cfg.Postgres.DBName),
		postgres.WithUsername(cfg.Postgres.User),
		postgres.WithPassword(cfg.Postgres.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(containerLogOccurrences).
				WithStartupTimeout(containerStartupTimeout),
		),
	)
	s.Require().NoError(err, "failed to start postgres container")

	cfg.Postgres.Host, err = s.container.Host(ctx)
	s.Require().NoError(err, "failed to get container host")

	port, err := s.container.MappedPort(ctx, "5432")
	s.Require().NoError(err, "failed to get container port")

	cfg.Postgres.Port = port.Port()

	err = migrations.RunMigrations(cfg.Postgres.URL())
	s.Require().NoError(err, "failed to run migrations")

	s.pool, err = New(ctx, &cfg.Postgres)
	s.Require().NoError(err, "failed to create pool")

	s.repo = NewRepository(s.pool)
}

func (s *RepositoryTestSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}

	if s.container != nil {
		err := testcontainers.TerminateContainer(s.container)
		s.Require().NoError(err, "failed to terminate container")
	}
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.cleanupTestData()
}

func (s *RepositoryTestSuite) cleanupTestData() {
	ctx := context.Background()

	query := `
        DO $$ 
        DECLARE 
            r RECORD;
        BEGIN
            FOR r IN (
                SELECT table_name 
                FROM information_schema.tables 
                WHERE table_schema = 'public' 
                  AND table_type = 'BASE TABLE'
                  AND table_name != 'schema_migrations'
            ) LOOP
                EXECUTE format('TRUNCATE TABLE %I RESTART IDENTITY CASCADE', r.table_name);
            END LOOP;
        END $$;
    `

	if _, err := s.pool.Exec(ctx, query); err != nil {
		s.T().Logf("warning: failed to cleanup test data: %v", err)
	}
}

//nolint:paralleltest // Integration tests share DB state
func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) assertAuthUserEqual(expected, actual model.AuthUser) {
	s.T().Helper()

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Login, actual.Login)
	s.Equal(expected.PasswordHash, actual.PasswordHash)
	s.Equal(expected.DisplayName, actual.DisplayName)
	s.Equal(expected.CreatedAt, actual.CreatedAt)
	s.Equal(expected.UpdatedAt, actual.UpdatedAt)
	s.True(actual.DeletedAt.IsZero())
}

func (s *RepositoryTestSuite) assertPgError(err error, expectedCode string) {
	s.T().Helper()

	s.Require().Error(err)

	var pgErr *pgconn.PgError
	s.Require().ErrorAs(err, &pgErr)
	s.Equal(expectedCode, pgErr.Code)
}

func (s *RepositoryTestSuite) assertUserEqual(expected, actual model.User) {
	s.T().Helper()

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Login, actual.Login)
	s.Equal(expected.DisplayName, actual.DisplayName)
	s.Equal(expected.CreatedAt, actual.CreatedAt)
	s.Equal(expected.UpdatedAt, actual.UpdatedAt)
	s.True(actual.DeletedAt.IsZero())
}

func (s *RepositoryTestSuite) generateTokenHash(prefix string, index int) string {
	hashSuffix := strconv.Itoa(index)
	padding := 64 - len(prefix) - len(hashSuffix)

	return prefix + hashSuffix + fmt.Sprintf("%0*d", padding, 0)
}

func (s *RepositoryTestSuite) createTestUser(ctx context.Context, loginSuffix string) model.AuthUser {
	s.T().Helper()

	user := model.AuthUser{
		User: model.User{
			Login:       loginSuffix,
			DisplayName: "TestUser",
		},
		PasswordHash: "hash_" + loginSuffix,
	}

	insertedUser, err := s.repo.InsertUser(ctx, user)
	s.Require().NoError(err)

	return insertedUser
}

func (s *RepositoryTestSuite) createTestChat(ctx context.Context) model.Chat {
	s.T().Helper()

	chat, err := s.repo.InsertChat(ctx, model.ChatTypePrivate)
	s.Require().NoError(err)

	return chat
}

func (s *RepositoryTestSuite) createTestToken(
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) model.Token {
	s.T().Helper()

	token := model.Token{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	insertedToken, err := s.repo.InsertRefreshToken(ctx, token)
	s.Require().NoError(err)

	return insertedToken
}
